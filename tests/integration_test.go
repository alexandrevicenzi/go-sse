package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexandrevicenzi/go-sse"
)

func TestServerShutdown(t *testing.T) {
	sseServer := sse.NewServer(nil)
	s := httptest.NewServer(sseServer)

	ch := make(chan struct{})
	ch2 := make(chan struct{})
	defer close(ch)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.URL, nil)
		if err != nil {
			t.Fail()
		}
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Connection", "keep-alive")

		_, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Fail()
		}

		close(ch2)
		// Wait until the main test routine is done
		<-ch
	}()
	// Wait until the HTTP request is sent
	<-ch2

	sseServer.Shutdown()
	s.Close()

	// Wait a bit after closing HTTP server
	time.Sleep(250 * time.Millisecond)
}
