package sse

import (
	"bytes"
	"fmt"
)

// Message represents a event source message.
type Message struct {
	id,
	data,
	event string
	retry int
}

func SimpleMessage(data string) *Message {
	return NewMessage("", data, "")
}

func NewMessage(id, data, event string) *Message {
	return &Message{
		id,
		data,
		event,
		0,
	}
}

func (m *Message) String() string {
	var buffer bytes.Buffer

	if len(m.id) > 0 {
		buffer.WriteString(fmt.Sprintf("id: %s\n", m.id))
	}

	if m.retry > 0 {
		buffer.WriteString(fmt.Sprintf("retry: %d\n", m.retry))
	}

	if len(m.event) > 0 {
		buffer.WriteString(fmt.Sprintf("event: %s\n", m.event))
	}

	if len(m.data) > 0 {
		buffer.WriteString(fmt.Sprintf("data: %s\n", m.data))
	}

	buffer.WriteString("\n")

	return buffer.String()
}
