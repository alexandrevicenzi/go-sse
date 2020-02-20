package main

import (
	"log"
	"strconv"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/gin-gonic/gin"
)

func main() {
	s := sse.NewServer(nil)
	defer s.Shutdown()

	router := gin.Default()

	router.StaticFile("/", "./static/index.html")

	router.GET("/events/:channel", func(c *gin.Context) {
		s.ServeHTTP(c.Writer, c.Request)
	})

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
	router.Run(":3000")
}
