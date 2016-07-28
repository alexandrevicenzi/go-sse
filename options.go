package sse

import (
    "net/http"
)

type Options struct {
    // RetryInterval change EventSource default retry interval (milliseconds).
    RetryInterval int
    // Headers allow to set custom headers (useful for CORS support).
    Headers map[string]string
    // ChannelNameFunc allow to create custom channel names.
    // Default channel name is the request path.
    ChannelNameFunc func (*http.Request) string
}

func (opt *Options) HasHeaders() bool {
    return opt.Headers != nil && len(opt.Headers) > 0
}
