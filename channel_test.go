package sse

import (
	"fmt"
	"sync"
	"testing"
)

func unused(i interface{}) {}

func readMsg(c *Client) {
	for msg := range c.send {
		unused(msg)
	}
}

func TestSendMessage(t *testing.T) {
	ch := newChannel("channel")
	defer ch.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	lastID := make(chan string, 1)
	msgCount := make(chan int, 1)

	c := newClient("", "client")

	go func() {
		i := 0
		id := ""

		for msg := range c.send {
			i++
			id = msg.id
		}

		lastID <- id
		msgCount <- i
	}()

	ch.addClient(c)

	go func() {
		for i := 0; i < 100; i++ {
			ch.SendMessage(NewMessage(fmt.Sprintf("id_%d", i+1), "msg", "channel"))
		}

		wg.Done()
	}()

	wg.Wait()

	ch.removeClient(c)

	if 100 != <-msgCount {
		t.Fatal("Wrong message count.")
	}

	if ch.LastEventID() != <-lastID {
		t.Fatal("Wrong Last ID.")
	}
}
