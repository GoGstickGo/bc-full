package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	bitcoinclient "node-dashboard/internal/client"
	"node-dashboard/internal/node"
	"node-dashboard/internal/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	// Command line flags for configuration
	rpcHost := flag.String("rpchost", "localhost:8332", "Bitcoin RPC host:port")
	rpcUser := flag.String("rpcuser", "", "Bitcoin RPC username")
	rpcPass := flag.String("rpcpass", "", "Bitcoin RPC password")
	useHTTPS := flag.Bool("https", false, "Use HTTPS for RPC connection")
	listenAddr := flag.String("listen", ":8080", "HTTP server listen address")
	updateInterval := flag.Duration("interval", 10*time.Second, "Stats update interval")
	flag.Parse()

	// Check required flags
	if *rpcUser == "" || *rpcPass == "" {
		log.Fatal("Bitcoin RPC username and password are required")
	}

	// Create Bitcoin RPC client
	client, err := bitcoinclient.NewClient(bitcoinclient.Config{
		Host:     *rpcHost,
		User:     *rpcUser,
		Pass:     *rpcPass,
		UseHTTPS: *useHTTPS,
	})
	if err != nil {
		log.Fatalf("Failed to create Bitcoin client: %v", err)
	}
	defer client.Shutdown()

	// Create WebSocket hub
	wsHub := websocket.NewHub()

	// Create node service
	service := node.NewService(client, wsHub)
	fmt.Printf("Service: %v\n", wsHub.GetLatestStats().BlockHeight)

	// Start data collection in a goroutine
	go service.Start(*updateInterval)
	defer service.Stop()

	// Set up Gin router
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Serve the main HTML page
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// API endpoint for current stats
	r.GET("/api/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, wsHub.GetLatestStats())
	})

	// Websocket endpoint
	r.GET("/ws", wsHub.HandleWebsocket)

	// Run server
	log.Printf("Starting server on %s", *listenAddr)
	if err := r.Run(*listenAddr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
