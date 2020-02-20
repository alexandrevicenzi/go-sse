package main

import (
	"strconv"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	"github.com/labstack/echo/v4"
)

func main() {
	s := sse.NewServer(nil)
	defer s.Shutdown()

	e := echo.New()
	e.File("/", "./static/index.html")

	e.Any("/events/:channel", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		s.ServeHTTP(res, req)
		return nil
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

	e.Logger.Fatal(e.Start(":3000"))
}
