package models

// NodeStats represents Bitcoin node statistics
type NodeStats struct {
	BlockHeight     int32                  `json:"blockHeight"`
	Connections     int32                  `json:"connections"`
	Difficulty      float64                `json:"difficulty"`
	MemPoolSize     int                    `json:"memPoolSize"`
	MemPoolBytes    int32                  `json:"memPoolBytes"`
	NetworkHashrate float64                `json:"networkHashrate"`
	MempoolFeeInfo  map[string]interface{} `json:"mempoolFeeInfo"`
	ChainInfo       interface{}            `json:"chainInfo"`
	LastUpdated     string                 `json:"lastUpdated"`
	Warnings        string                 `json:"warnings"`
}
