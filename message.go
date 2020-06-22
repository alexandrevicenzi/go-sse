package sse

import (
	"bytes"
	"fmt"
	"strings"
)

// Message represents a event source message.
type Message struct {
	id,
	data,
	event string
	retry int
}

// SimpleMessage creates a simple event source message.
func SimpleMessage(data string) *Message {
	return NewMessage("", data, "")
}

// NewMessage creates an event source message.
func NewMessage(id, data, event string) *Message {
	return &Message{
		id,
		data,
		event,
		0,
	}
}

func (m *Message) Buffer() *bytes.Buffer {
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
		buffer.WriteString(fmt.Sprintf("data: %s\n", strings.Replace(m.data, "\n", "\ndata: ", -1)))
	}

	buffer.WriteString("\n")

	return &buffer
}

func (m *Message) String() string {
	return m.Buffer().String()
}

func (m *Message) Bytes() []byte {
	return m.Buffer().Bytes()
}
