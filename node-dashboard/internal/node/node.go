package node

import (
	"log"
	"time"

	bitcoinclient "node-dashboard/internal/client"
	"node-dashboard/internal/models"
	"node-dashboard/internal/websocket"
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
	}

	// Get mempool info using GetRawMempoolVerbose
	mempoolInfo, err := s.client.GetRawMempoolVerbose()
	if err != nil {
		log.Printf("Error getting mempool info: %v", err)
	} else {
		// Count the number of transactions in the mempool
		nodeStats.MemPoolSize = len(mempoolInfo)

		// Calculate total size of all transactions
		var totalSize int32
		for _, tx := range mempoolInfo {
			totalSize += tx.Size
		}
		nodeStats.MemPoolBytes = totalSize
	}

	// Get network hashrate (using 120 blocks)
	hashrate, err := s.client.GetNetworkHashPS()
	if err != nil {
		log.Printf("Error getting network hashrate: %v", err)
	} else {
		nodeStats.NetworkHashrate = hashrate
	}

	// Get mempool fee info - Using both ECONOMICAL and CONSERVATIVE modes
	/*economicalMode := btcjson.EstimateModeEconomical
	conservativeMode := btcjson.EstimateModeConservative

	// Get economical fee estimate
	economicalFeeInfo, err := s.client.EstimateSmartFee(2, &economicalMode)
	if err != nil {
		log.Printf("Error getting economical fee estimate: %v", err)
	}

	// Get conservative fee estimate
	conservativeFeeInfo, err := s.client.EstimateSmartFee(2, &conservativeMode)
	if err != nil {
		log.Printf("Error getting conservative fee estimate: %v", err)
	}

	// Store both estimates
	nodeStats.MempoolFeeInfo = map[string]interface{}{
		"economical":   economicalFeeInfo,
		"conservative": conservativeFeeInfo,
	}*/

	nodeStats.LastUpdated = time.Now().Format(time.RFC3339)

	// Update the websocket hub with the new data
	s.wsHub.UpdateStats(nodeStats)
}
