package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/gorilla/mux"
)

func main() {
	s := sse.NewServer(nil)
	defer s.Shutdown()

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.Handle("/events/{channel}", s)

	go func() {
		for {
			s.SendMessage("/events/channel-1", sse.SimpleMessage(time.Now().Format("2006/02/01/ 15:04:05")))
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		i := 0
		for {
			i++
			s.SendMessage("/events/channel-2", sse.SimpleMessage(strconv.Itoa(i)))
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Listening at :3000")
	http.ListenAndServe(":3000", r)
}
