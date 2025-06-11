package models

// NodeStats represents Bitcoin node statistics
type NodeStats struct {
	BlockHeight     int32       `json:"blockHeight"`
	Connections     int32       `json:"connections"`
	Difficulty      float64     `json:"difficulty"`
	NetworkHashrate float64     `json:"networkHashrate"`
	ChainInfo       interface{} `json:"chainInfo"`
	LastUpdated     string      `json:"lastUpdated"`
}
