package sse

// Channel represents a server sent events channel.
type Channel struct {
	lastEventID,
	name string
	clients map[*Client]bool
}

func newChannel(name string) *Channel {
	return &Channel{
		"",
		name,
		make(map[*Client]bool),
	}
}

// SendMessage broadcast a message to all clients in a channel.
func (c *Channel) SendMessage(message *Message) {
	c.lastEventID = message.id

	for c, open := range c.clients {
		if open {
			c.send <- message
		}
	}
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
	return len(c.clients)
}

// LastEventID returns the ID of the last message sent.
func (c *Channel) LastEventID() string {
	return c.lastEventID
}

func (c *Channel) addClient(client *Client) {
	c.clients[client] = true
}

func (c *Channel) removeClient(client *Client) {
	c.clients[client] = false
	close(client.send)
	delete(c.clients, client)
}
