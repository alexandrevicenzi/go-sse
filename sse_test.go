package sse

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewServerNilOptions(t *testing.T) {
	srv := NewServer(nil)
	defer srv.Shutdown()

	if srv == nil || srv.options == nil || srv.options.Logger == nil {
		t.Fail()
	}
}

func TestNewServerNilLogger(t *testing.T) {
	srv := NewServer(&Options{
		Logger: nil,
	})

	defer srv.Shutdown()

	if srv == nil || srv.options == nil || srv.options.Logger == nil {
		t.Fail()
	}
}

func TestServer(t *testing.T) {
	srv := NewServer(&Options{
		Logger: nil,
	})

	defer srv.Shutdown()

	for n := 0; n < 2; n++ {
		name := fmt.Sprintf("CH-%d", n+1)
		srv.addChannel(name)
		// fmt.Printf("Channel %s registed\n", name)
	}

	messages := make([]*Message, 0)
	wgReg := sync.WaitGroup{}
	wgCh := sync.WaitGroup{}

	for n := 0; n < 5; n++ {
		for name, ch := range srv.channels {
			wgReg.Add(1)
			wgCh.Add(1)

			go func(id string) {
				c := newClient("", name)
				ch.addClient(c)

				// fmt.Printf("Client %s registed to channel %s\n", id, name)
				wgReg.Done()

				for msg := range c.send {
					// fmt.Printf("Channel: %s - Client: %s - Message: %s\n", name, id, msg.data)
					messages = append(messages, msg)
				}

				wgCh.Done()
			}(fmt.Sprintf("C-%d", n+1))
		}
	}

	wgReg.Wait()

	srv.SendMessage("", SimpleMessage("hello"))

	srv.close()

	wgCh.Wait()

	if len(messages) != 10 {
		t.Fail()
	}
}
