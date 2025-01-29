package internal

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Client struct {
	clientsMutex sync.Mutex
	clients      []chan []byte
}

func NewClient() *Client {
	return &Client{}
}

// Add client to the list
func (c *Client) AddClient() chan []byte {
	client := make(chan []byte, 2048)
	c.clientsMutex.Lock()
	c.clients = append(c.clients, client)
	c.clientsMutex.Unlock()
	return client
}

// Remove client from the list
func (c *Client) RemoveClient(client chan []byte) {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	// for i, c := range c.clients {
	// 	if c == client {
	// 		c.clients = append(clients[:i], clients[i+1:]...)
	// 		close(c)
	// 		break
	// 	}
	// }
	newClients := make([]chan []byte, 0, len(c.clients)-1)
	for _, c := range c.clients {
		if c == client {
			close(c)
		} else {
			newClients = append(newClients, c)
		}
	}
	c.clients = newClients
}

// Broadcast data to all clients
func (c *Client) Broadcast(data []byte) {
	c.clientsMutex.Lock()
	defer c.clientsMutex.Unlock()

	for _, client := range c.clients {
		select {
		case client <- data:
		default:
			c.RemoveClient(client)
		}
	}
}

// Get all clients
func (c *Client) GetClients() []chan []byte {
	return c.clients
}

// HandleClient handles the client
func (c *Client) HandleClient(w http.ResponseWriter, _ *http.Request) {
	client := c.AddClient()
	defer c.RemoveClient(client)
	w.Header().Set("Content-Type", "audio/mpeg")
	// w.Header().Set("Transfer-Encoding", "chunked") // Enable chunked transfer encoding
	// w.Header().Set("Connection", "keep-alive") // Keep the connection open

	for data := range client {
		if _, err := w.Write(data); err != nil {
			log.Println("Error writing to client:", err)
			break
		}
	}
}

// HandleClientGin handles the client for gin library
func (c *Client) HandleClientGin(g *gin.Context) {
	g.Header("Content-Type", "audio/mpeg")
	client := c.AddClient()
	defer c.RemoveClient(client)
	for data := range client {
		if _, err := g.Writer.Write(data); err != nil {
			log.Println("Error writing to client:", err)
			break
		}
	}
}
