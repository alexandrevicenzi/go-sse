package sse

import (
	"context"
	"fmt"
	"go.uber.org/goleak"
	"io/ioutil"
	"log"
	"net/http"
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

	// Create N channels
	for n := 0; n < channelCount; n++ {
		name := fmt.Sprintf("CH-%d", n+1)
		srv.addChannel(name)
		fmt.Printf("Channel %s registed\n", name)
	}

	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	// Create N clients in all channels
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

func TestShutdown(t *testing.T) {
	defer goleak.VerifyNone(t)

	srv := NewServer(nil)

	http.Handle("/events/", srv)

	httpServer := &http.Server{Addr: ":3000", Handler: nil}

	go func() { _ = httpServer.ListenAndServe() }()

	stop := make(chan struct{})

	go func() {
		r, err := http.Get("http://localhost:3000/events/chan")
		if err != nil {
			log.Fatalln(err)
			return
		}
		// Stop while client is reading the response
		stop <- struct{}{}
		_, _ = ioutil.ReadAll(r.Body)
	}()

	<-stop

	srv.Shutdown()

	ctx, done := context.WithTimeout(context.Background(), 600*time.Millisecond)
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	done()
}
