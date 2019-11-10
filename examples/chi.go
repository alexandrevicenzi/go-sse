package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.DefaultCompress)

	s := sse.NewServer(nil)
	defer s.Shutdown()

	r.Handle("/", http.FileServer(http.Dir("./static")))
	r.Mount("/events/", s)

	go func() {
		for {
			s.SendMessage("/events/channel-1", sse.SimpleMessage(time.Now().String()))
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
