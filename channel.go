package sse

import (
	"sync"
)

// Channel represents a server sent events channel.
type Channel struct {
	mu          sync.RWMutex
	lastEventID string
	name        string
	clients     map[*Client]bool
}

func newChannel(name string) *Channel {
	return &Channel{
		sync.RWMutex{},
		"",
		name,
		make(map[*Client]bool),
	}
}

// SendMessage broadcast a message to all clients in a channel.
func (c *Channel) SendMessage(message *Message) {
	c.lastEventID = message.id

	c.mu.RLock()

	for c, open := range c.clients {
		if open {
			c.send <- message
		}
	}

	c.mu.RUnlock()
}

// Close closes the channel and disconnect all clients.
func (c *Channel) Close() {
	// Kick all clients of this channel.
	for client := range c.clients {
		c.removeClient(client)
	}
}

// ClientCount returns the number of clients connected to this channel.
func (c *Channel) ClientCount() int {
	c.mu.RLock()
	count := len(c.clients)
	c.mu.RUnlock()

	return count
}

// LastEventID returns the ID of the last message sent.
func (c *Channel) LastEventID() string {
	return c.lastEventID
}

func (c *Channel) addClient(client *Client) {
	c.mu.Lock()
	c.clients[client] = true
	c.mu.Unlock()
}

func (c *Channel) removeClient(client *Client) {
	c.mu.Lock()
	c.clients[client] = false
	delete(c.clients, client)
	c.mu.Unlock()
	close(client.send)
}
