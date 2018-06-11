# go-sse [![Go Report Card](https://goreportcard.com/badge/github.com/alexandrevicenzi/go-sse)](https://goreportcard.com/report/github.com/alexandrevicenzi/go-sse) [![Build Status](https://travis-ci.org/alexandrevicenzi/go-sse.svg?branch=master)](https://travis-ci.org/alexandrevicenzi/go-sse) [![GoDoc](https://godoc.org/github.com/alexandrevicenzi/go-sse?status.svg)](http://godoc.org/github.com/alexandrevicenzi/go-sse)

Server-Sent Events for Go

## About

[Server-sent events](http://www.html5rocks.com/en/tutorials/eventsource/basics/) is a method of continuously sending data from a server to the browser, rather than repeatedly requesting it, replacing the "long polling way".

It's [supported](http://caniuse.com/#feat=eventsource) by all major browsers and for IE/Edge you can use a [polyfill](https://github.com/Yaffle/EventSource).

`go-sse` is a small library to create a Server-Sent Events server in Go.

## Features

- Multiple channels (isolated)
- Broadcast message to all channels
- Custom headers (useful for CORS)
- `Last-Event-ID` support (resend lost messages)
- [Follow SSE specification](https://html.spec.whatwg.org/multipage/comms.html#server-sent-events)

## Getting Started

Simple Go example that handle 2 channels and send messages to all clients connected in each channel.

```go
package main

import (
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/alexandrevicenzi/go-sse"
)

func main() {
    // Create the server.
    s := sse.NewServer(nil)
    defer s.Shutdown()

    // Register with /events endpoint.
    http.Handle("/events/", s)

    // Dispatch messages to channel-1.
    go func () {
        for {
            s.SendMessage("/events/channel-1", sse.SimpleMessage(time.Now().String()))
            time.Sleep(5 * time.Second)
        }
    }()

    // Dispatch messages to channel-2
    go func () {
        i := 0
        for {
            i++
            s.SendMessage("/events/channel-2", sse.SimpleMessage(strconv.Itoa(i)))
            time.Sleep(5 * time.Second)
        }
    }()

    http.ListenAndServe(":3000", nil)
}
```

Connecting to our server from JavaScript:

```js
e1 = new EventSource('/events/channel-1');
e1.onmessage = function(event) {
    // do something...
};

e2 = new EventSource('/events/channel-2');
e2.onmessage = function(event) {
    // do something...
};
```
