package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"node-dashboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Hub manages WebSocket connections
type Hub struct {
	clients     map[*websocket.Conn]bool
	clientsMu   sync.Mutex
	latestStats models.NodeStats
	statsMu     sync.RWMutex
	upgrader    websocket.Upgrader
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections (customize for production)
			},
		},
	}
}

// HandleWebsocket upgrades HTTP connection to WebSocket and manages the connection
func (h *Hub) HandleWebsocket(c *gin.Context) {
	ws, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	// Add new client to the map
	h.clientsMu.Lock()
	h.clients[ws] = true
	h.clientsMu.Unlock()

	// Send current stats immediately
	h.statsMu.RLock()
	data, err := json.Marshal(h.latestStats)
	h.statsMu.RUnlock()

	if err != nil {
		log.Printf("Error writing message: %v", err)
		return
	}
	ws.WriteMessage(websocket.TextMessage, data)

	// Handle disconnection
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			h.clientsMu.Lock()
			delete(h.clients, ws)
			h.clientsMu.Unlock()
			break
		}
	}
}

// UpdateStats updates the cached stats and broadcasts to all clients
func (h *Hub) UpdateStats(stats models.NodeStats) {
	// Update the cached stats
	h.statsMu.Lock()
	h.latestStats = stats
	h.statsMu.Unlock()

	// Broadcast to all clients
	h.Broadcast()
}

// Broadcast sends the latest stats to all connected websocket clients
func (h *Hub) Broadcast() {
	h.statsMu.RLock()
	data, err := json.Marshal(h.latestStats)
	h.statsMu.RUnlock()

	if err != nil {
		log.Printf("Error marshaling stats: %v", err)
		return
	}

	h.clientsMu.Lock()
	for client := range h.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error sending to client: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
	h.clientsMu.Unlock()
}

// GetLatestStats returns the most recent node stats
func (h *Hub) GetLatestStats() models.NodeStats {
	h.statsMu.RLock()
	defer h.statsMu.RUnlock()
	if h.latestStats.BlockHeight == 0 {
		log.Print("No stats available")
		return models.NodeStats{}
	}
	return h.latestStats
}
