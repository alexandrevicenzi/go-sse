package sse

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"
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
	channelCount := 2
	clientCount := 5
	messageCount := 0

	srv := NewServer(&Options{
		Logger: log.New(os.Stdout, "go-sse: ", log.Ldate|log.Ltime|log.Lshortfile),
	})

	defer srv.Shutdown()

	// Create N channes
	for n := 0; n < channelCount; n++ {
		name := fmt.Sprintf("CH-%d", n+1)
		srv.addChannel(name)
		fmt.Printf("Channel %s registed\n", name)
	}

	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	// Create N clients in all channes
	for n := 0; n < clientCount; n++ {
		for name, ch := range srv.channels {
			wg.Add(1)

			// Create new client
			c := newClient("", name)
			// Add client to current channel
			ch.addClient(c)

			id := fmt.Sprintf("C-%d", n+1)
			fmt.Printf("Client %s registed to channel %s\n", id, name)

			go func(id string) {
				// Wait for messages in the channel
				for msg := range c.send {
					m.Lock()
					messageCount++
					m.Unlock()
					fmt.Printf("Channel: %s - Client: %s - Message: %s\n", name, id, msg.data)
					wg.Done()
				}
			}(id)
		}
	}

	// Send hello message to all channels and all clients in it
	srv.SendMessage("", SimpleMessage("hello"))

	srv.close()

	wg.Wait()

	if messageCount != channelCount*clientCount {
		t.Errorf("Expected %d messages but got %d", channelCount*clientCount, messageCount)
	}
}

func TestMultipleTopics(t *testing.T) {
	// usage pattern we have a client which subsribed to multiple channels
	// in one connection
	sendersWg := sync.WaitGroup{}
	workerWg := sync.WaitGroup{}
	m := sync.Mutex{}
	messageCount := 0
	numMessages := 3
	topics := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thuesday", "Friday", "Saturday"}

	srv := NewServer(&Options{
		Logger: log.New(os.Stdout, "go-sse: ", log.Ldate|log.Ltime|log.Lshortfile),
	})

	defer srv.Shutdown()

	name := "CH-0"
	ch := srv.addChannel(name)
	fmt.Printf("Channel %s registed\n", name)

	// Create new client
	c := newClient("", name)
	// Add client to current channel
	ch.addClient(c)

	// receive messages
	workerWg.Add(1)
	go func() {
		defer workerWg.Done()
		// Wait for messages in the channel
		for msg := range c.send {
			m.Lock()
			messageCount++
			m.Unlock()
			fmt.Printf("ID: %s - Topic: %s - Message: %s\n", msg.id, msg.event, msg.data)
		}
	}()

	for _, day := range topics {
		sendersWg.Add(1)
		go func(topic string) {
			defer sendersWg.Done()
			for n := 0; n < numMessages; n++ {
				srv.SendMessage(name, NewMessage(id(), "hello", topic))
			}
		}(day)
	}
	// Wait senders to complete
	sendersWg.Wait()

	srv.close()
	// Wait recipient routine
	workerWg.Wait()

	if messageCount != len(topics)*numMessages {
		t.Errorf("Expected %d messages but got %d", 3*numMessages, messageCount)
	}

}

func id() string {
	return fmt.Sprintf("%d", time.Now().UTC().Unix()*1000)
}
