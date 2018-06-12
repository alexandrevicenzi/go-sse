package sse

import (
	"log"
	"net/http"
)

// Options holds server configurations.
type Options struct {
	// RetryInterval change EventSource default retry interval (milliseconds).
	RetryInterval int
	// Headers allow to set custom headers (useful for CORS support).
	Headers map[string]string
	// ChannelNameFunc allow to create custom channel names.
	// Default channel name is the request path.
	ChannelNameFunc func(*http.Request) string
	// All usage logs end up in Logger
	Logger *log.Logger
}

func (opt *Options) hasHeaders() bool {
	return opt.Headers != nil && len(opt.Headers) > 0
}
