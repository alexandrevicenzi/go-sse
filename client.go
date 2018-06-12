package sse

// Client represents a web browser connection.
type Client struct {
	lastEventID,
	channel string
	send chan *Message
}

func newClient(lastEventID, channel string) *Client {
	return &Client{
		lastEventID,
		channel,
		make(chan *Message),
	}
}

// SendMessage sends a message to client.
func (c *Client) SendMessage(message *Message) {
	c.lastEventID = message.id
	c.send <- message
}

// Channel returns the channel where this client is subscribe to.
func (c *Client) Channel() string {
	return c.channel
}

// LastEventID returns the ID of the last message sent.
func (c *Client) LastEventID() string {
	return c.lastEventID
}
