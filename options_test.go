package sse

import "testing"

func TestHasHeadersEmpty(t *testing.T) {
	opt := Options{}

	if opt.hasHeaders() == true {
		t.Fatal("There's headers.")
	}
}

func TestHasHeadersNotEmpty(t *testing.T) {
	opt := Options{}

	opt = Options{
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
		},
	}

	if opt.hasHeaders() == false {
		t.Fatal("There's no headers.")
	}
}
