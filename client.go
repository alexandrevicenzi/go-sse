package sse

type Client struct {
    lastEventId,
    channel string
    send chan *Message
}

func NewClient(lastEventId, channel string) *Client {
    return &Client{
        lastEventId,
        channel,
        make(chan *Message),
    }
}

// SendMessage sends a message to client.
func (c *Client) SendMessage(message *Message) {
    c.lastEventId = message.id
    c.send <- message
}

// Channel returns the channel where this client is subscribe to.
func (c *Client) Channel() string {
    return c.channel
}

func (c *Client) LastEventId() string {
    return c.lastEventId
}
