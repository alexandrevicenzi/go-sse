package main

import (
    "log"
    "net/http"
    "strconv"
    "time"
    "os"

    "github.com/alexandrevicenzi/go-sse"
)

func main() {
    s := sse.NewServer(&sse.Options{
        // Increase default retry interval to 10s.
        RetryInterval: 10 * 1000,
        // CORS headers
        Headers: map[string]string {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "GET, OPTIONS",
            "Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
        },
        // Custom channel name generator
        ChannelNameFunc: func (request *http.Request) string {
            return request.URL.Path
        },
        // Print debug info
        Logger: log.New(os.Stdout,
          "go-sse: ",
          log.Ldate|log.Ltime|log.Lshortfile),
    })

    defer s.Shutdown()

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.Handle("/events/", s)

    go func () {
        for {
            s.SendMessage("/events/channel-1", sse.SimpleMessage(time.Now().String()))
            time.Sleep(5 * time.Second)
        }
    }()

    go func () {
        i := 0
        for {
            i++
            s.SendMessage("/events/channel-2", sse.SimpleMessage(strconv.Itoa(i)))
            time.Sleep(5 * time.Second)
        }
    }()

    log.Println("Listening at :3000")
    http.ListenAndServe(":3000", nil)
}
