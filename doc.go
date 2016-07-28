// Package sse implements Server-Sent Events that supports multiple channels.
//
// Server-sent events is a method of continuously sending data from a server to the browser, rather than repeatedly requesting it.
//
// Examples
//
// Basic usage of sse package.
//
//    s := sse.NewServer(nil)
//    defer s.Shutdown()
//
//    http.Handle("/events/", s)
//
package sse
