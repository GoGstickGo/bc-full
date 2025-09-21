package node

import (
	"log"
	"strings"
	"time"

	bitcoinclient "node-dashboard/internal/client"
	"node-dashboard/internal/models"
	"node-dashboard/internal/websocket"

	"github.com/btcsuite/btcd/btcjson"
)

// Service handles Bitcoin node data collection and processing
type Service struct {
	client *bitcoinclient.Client
	wsHub  *websocket.Hub
	stop   chan struct{}
}

// NewService creates a new node service
func NewService(client *bitcoinclient.Client, wsHub *websocket.Hub) *Service {
	return &Service{
		client: client,
		wsHub:  wsHub,
		stop:   make(chan struct{}),
	}
}

// Start begins the node data collection process
func (s *Service) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial data fetch
	s.fetchNodeStats()

	// Continuous data fetch
	for {
		select {
		case <-ticker.C:
			s.fetchNodeStats()
		case <-s.stop:
			log.Println("Node service stopping")
			return
		}
	}
}

// Stop halts the data collection process
func (s *Service) Stop() {
	close(s.stop)
}

// fetchNodeStats collects data from the Bitcoin node
func (s *Service) fetchNodeStats() {
	var nodeStats models.NodeStats

	// Get blockchain info
	blockchainInfo, err := s.client.GetBlockChainInfo()
	if err != nil {
		log.Printf("Error getting blockchain info: %v", err)
	} else {
		nodeStats.BlockHeight = blockchainInfo.Blocks
		nodeStats.Difficulty = blockchainInfo.Difficulty
		nodeStats.ChainInfo = blockchainInfo
	}

	// Get network info
	networkInfo, err := s.client.GetNetworkInfo()
	if err != nil {
		log.Printf("Error getting network info: %v", err)
	} else {
		nodeStats.Connections = networkInfo.Connections
		nodeStats.Version = networkInfo.Version

		onionReachable := false
		if networkInfo.Networks != nil {
			for _, network := range networkInfo.Networks {
				if network.Name == "onion" && network.Reachable {
					onionReachable = true
					break
				}
			}
		}

		nodeStats.OnionStatus = onionReachable
	}

	// Get network hashrate (using 120 blocks)
	hashrate, err := s.client.GetNetworkHashPS()
	if err != nil {
		log.Printf("Error getting network hashrate: %v", err)
	} else {
		nodeStats.NetworkHashrate = hashrate
	}

	balance, err := s.client.GetBalance("*")
	if err != nil {
		log.Printf("Error getting node info: %v", err)
	} else {
		nodeStats.Balance = balance.ToBTC()
	}

	// Get mempool info using GetRawMempoolVerbose
	mempoolInfo, err := s.client.GetRawMempoolVerbose()
	if err != nil {
		log.Printf("Error getting mempool info: %v", err)
	} else {
		// Count the number of transactions in the mempool
		nodeStats.MemPoolSize = len(mempoolInfo)

		// Calculate total size of all transactions
		var totalSize int64
		for _, tx := range mempoolInfo {
			totalSize += int64(tx.Vsize)
		}
		nodeStats.MemPoolBytes = totalSize
	}

	feeEstimates := make(map[string]interface{})

	// Get economical fee estimate
	economicalFeeInfo, err := s.client.EstimateSmartFee(2, &btcjson.EstimateModeEconomical)
	if err != nil {
		if strings.Contains(err.Error(), "Fee estimation disabled") {
			log.Println("Fee estimation is disabled on this node")
			feeEstimates["economical"] = map[string]interface{}{
				"error":   "Fee estimation disabled",
				"feerate": nil,
			}
		} else {
			log.Printf("Error getting economical fee estimate: %v", err)
			feeEstimates["economical"] = map[string]interface{}{
				"error":   err.Error(),
				"feerate": nil,
			}
		}
	} else {
		feeEstimates["economical"] = economicalFeeInfo
	}

	// Get conservative fee estimate
	conservativeFeeInfo, err := s.client.EstimateSmartFee(2, &btcjson.EstimateModeConservative)
	if err != nil {
		if strings.Contains(err.Error(), "Fee estimation disabled") {
			log.Println("Conservative fee estimation is also disabled")
			feeEstimates["conservative"] = map[string]interface{}{
				"error":   "Fee estimation disabled",
				"feerate": nil,
			}
		} else {
			log.Printf("Error getting conservative fee estimate: %v", err)
			feeEstimates["conservative"] = map[string]interface{}{
				"error":   err.Error(),
				"feerate": nil,
			}
		}
	} else {
		feeEstimates["conservative"] = conservativeFeeInfo
	}

	nodeStats.MemPoolFeeInfo = feeEstimates

	nodeStats.LastUpdated = time.Now().Format(time.RFC3339)

	// Update the websocket hub with the new data
	s.wsHub.UpdateStats(nodeStats)
}
