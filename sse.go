package sse

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

// Server represents a server sent events server.
type Server struct {
	mu           sync.RWMutex
	options      *Options
	channels     map[string]*Channel
	addClient    chan *Client
	removeClient chan *Client
	shutdown     chan bool
	closeChannel chan string
}

// NewServer creates a new SSE server.
func NewServer(options *Options) *Server {
	if options == nil {
		options = &Options{
			Logger: log.New(os.Stdout, "go-sse: ", log.LstdFlags),
		}
	}

	if options.Logger == nil {
		options.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	s := &Server{
		sync.RWMutex{},
		options,
		make(map[string]*Channel),
		make(chan *Client),
		make(chan *Client),
		make(chan bool),
		make(chan string),
	}

	go s.dispatch()

	return s
}

func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	flusher, ok := response.(http.Flusher)

	if !ok {
		http.Error(response, "Streaming unsupported.", http.StatusInternalServerError)
		return
	}

	h := response.Header()

	if s.options.hasHeaders() {
		for k, v := range s.options.Headers {
			h.Set(k, v)
		}
	}

	if request.Method == "GET" {
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("X-Accel-Buffering", "no")

		var channelName string

		if s.options.ChannelNameFunc == nil {
			channelName = request.URL.Path
		} else {
			channelName = s.options.ChannelNameFunc(request)
		}

		lastEventID := request.Header.Get("Last-Event-ID")
		c := newClient(lastEventID, channelName)
		s.addClient <- c
		closeNotify := request.Context().Done()

		go func() {
			<-closeNotify
			s.removeClient <- c
		}()

		response.WriteHeader(http.StatusOK)
		flusher.Flush()

		for msg := range c.send {
			msg.retry = s.options.RetryInterval
			response.Write(msg.Bytes())
			flusher.Flush()
		}
	} else if request.Method != "OPTIONS" {
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SendMessage broadcast a message to all clients in a channel.
// If channelName is an empty string, it will broadcast the message to all channels.
func (s *Server) SendMessage(channelName string, message *Message) {
	if len(channelName) == 0 {
		s.options.Logger.Print("broadcasting message to all channels.")

		s.mu.RLock()

		for _, ch := range s.channels {
			ch.SendMessage(message)
		}

		s.mu.RUnlock()
	} else if ch, ok := s.getChannel(channelName); ok {
		ch.SendMessage(message)
		s.options.Logger.Printf("message sent to channel '%s'.", channelName)
	} else {
		s.options.Logger.Printf("message not sent because channel '%s' has no clients.", channelName)
	}
}

// Restart closes all channels and clients and allow new connections.
func (s *Server) Restart() {
	s.options.Logger.Print("restarting server.")
	s.close()
}

// Shutdown performs a graceful server shutdown.
func (s *Server) Shutdown() {
	s.shutdown <- true
}

// ClientCount returns the number of clients connected to this server.
func (s *Server) ClientCount() int {
	i := 0

	s.mu.RLock()

	for _, channel := range s.channels {
		i += channel.ClientCount()
	}

	s.mu.RUnlock()

	return i
}

// HasChannel returns true if the channel associated with name exists.
func (s *Server) HasChannel(name string) bool {
	_, ok := s.getChannel(name)
	return ok
}

// GetChannel returns the channel associated with name or nil if not found.
func (s *Server) GetChannel(name string) (*Channel, bool) {
	return s.getChannel(name)
}

// Channels returns a list of all channels to the server.
func (s *Server) Channels() []string {
	channels := []string{}

	s.mu.RLock()

	for name := range s.channels {
		channels = append(channels, name)
	}

	s.mu.RUnlock()

	return channels
}

// CloseChannel closes a channel.
func (s *Server) CloseChannel(name string) {
	s.closeChannel <- name
}

func (s *Server) addChannel(name string) *Channel {
	ch := newChannel(name)

	s.mu.Lock()
	s.channels[ch.name] = ch
	s.mu.Unlock()

	s.options.Logger.Printf("channel '%s' created.", ch.name)

	return ch
}

func (s *Server) removeChannel(ch *Channel) {
	s.mu.Lock()
	delete(s.channels, ch.name)
	s.mu.Unlock()

	ch.Close()

	s.options.Logger.Printf("channel '%s' closed.", ch.name)
}

func (s *Server) getChannel(name string) (*Channel, bool) {
	s.mu.RLock()
	ch, ok := s.channels[name]
	s.mu.RUnlock()
	return ch, ok
}

func (s *Server) close() {
	for _, ch := range s.channels {
		s.removeChannel(ch)
	}
}

func (s *Server) dispatch() {
	s.options.Logger.Print("server started.")

	for {
		select {

		// New client connected.
		case c := <-s.addClient:
			ch, exists := s.getChannel(c.channel)

			if !exists {
				ch = s.addChannel(c.channel)
			}

			ch.addClient(c)
			s.options.Logger.Printf("new client connected to channel '%s'.", ch.name)

		// Client disconnected.
		case c := <-s.removeClient:
			if ch, exists := s.getChannel(c.channel); exists {
				ch.removeClient(c)
				s.options.Logger.Printf("client disconnected from channel '%s'.", ch.name)

				if ch.ClientCount() == 0 {
					s.options.Logger.Printf("channel '%s' has no clients.", ch.name)
					s.removeChannel(ch)
				}
			}

		// Close channel and all clients in it.
		case channel := <-s.closeChannel:
			if ch, exists := s.getChannel(channel); exists {
				s.removeChannel(ch)
			} else {
				s.options.Logger.Printf("requested to close nonexistent channel '%s'.", channel)
			}

		// Event Source shutdown.
		case <-s.shutdown:
			s.close()
			close(s.addClient)
			close(s.removeClient)
			close(s.closeChannel)
			close(s.shutdown)

			s.options.Logger.Print("server stopped.")
			return
		}
	}
}
