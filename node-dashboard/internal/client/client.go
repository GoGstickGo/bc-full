package bitcoinclient

import (
	"log"

	"github.com/btcsuite/btcd/rpcclient"
)

// Client wraps the btcd RPC client
type Client struct {
	*rpcclient.Client
}

// Config holds Bitcoin RPC connection configuration
type Config struct {
	Host     string
	User     string
	Pass     string
	UseHTTPS bool
}

// NewClient creates a new Bitcoin RPC client
func NewClient(config Config) (*Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         config.Host,
		User:         config.User,
		Pass:         config.Pass,
		HTTPPostMode: true,
		DisableTLS:   !config.UseHTTPS,
	}

	// Create a new RPC client using the provided configuration
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Printf("Error creating Bitcoin RPC client: %v", err)
		return nil, err
	}

	return &Client{client}, nil
}

// Shutdown safely closes the connection to the Bitcoin node
func (c *Client) Shutdown() {
	c.Client.Shutdown()
	log.Println("Bitcoin RPC client shutdown complete")
}
