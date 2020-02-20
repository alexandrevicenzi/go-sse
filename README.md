# go-sse

[![Go Report Card](https://goreportcard.com/badge/github.com/alexandrevicenzi/go-sse)](https://goreportcard.com/report/github.com/alexandrevicenzi/go-sse)
[![Build Status](https://travis-ci.org/alexandrevicenzi/go-sse.svg?branch=master)](https://travis-ci.org/alexandrevicenzi/go-sse)
[![GoDoc](https://godoc.org/github.com/alexandrevicenzi/go-sse?status.svg)](http://godoc.org/github.com/alexandrevicenzi/go-sse)

Server-Sent Events for Go

## About

[Server-sent events](http://www.html5rocks.com/en/tutorials/eventsource/basics/) is a method of continuously sending data from a server to the browser, rather than repeatedly requesting it, replacing the "long polling way".

It's [supported](http://caniuse.com/#feat=eventsource) by all major browsers and for IE/Edge you can use a [polyfill](https://github.com/Yaffle/EventSource).

`go-sse` is a small library to create a Server-Sent Events server in Go and works with Go 1.9+.

## Features

- Multiple channels (isolated)
- Broadcast message to all channels
- Custom headers (useful for CORS)
- `Last-Event-ID` support (resend lost messages)
- [Follow SSE specification](https://html.spec.whatwg.org/multipage/comms.html#server-sent-events)
- Compatible with multiple Go frameworks

## Instalation

`go get github.com/alexandrevicenzi/go-sse`

## Example

Server side:

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
    // Create SSE server
    s := sse.NewServer(nil)
    defer s.Shutdown()

    // Configure the route
    http.Handle("/events/", s)

    // Send messages every 5 seconds
    go func() {
        for {
            s.SendMessage("/events/my-channel", sse.SimpleMessage(time.Now().Format("2006/02/01/ 15:04:05")))
            time.Sleep(5 * time.Second)
        }
    }()

    log.Println("Listening at :3000")
    http.ListenAndServe(":3000", nil)
}
```

Client side (JavaScript):

```js
e = new EventSource('/events/my-channel');
e.onmessage = function(event) {
    console.log(event.data);
};
```

More examples available [here](https://github.com/alexandrevicenzi/go-sse/tree/master/_examples).

## License

MIT
