package sse

import "testing"

func TestEmptyMessage(t *testing.T) {
	msg := Message{}

	if msg.String() != "\n" {
		t.Fatal("Messagem not empty.")
	}
}

func TestDataMessage(t *testing.T) {
	msg := Message{data: "test"}

	if msg.String() != "data: test\n\n" {
		t.Fatal("Messagem not empty.")
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		"123",
		"test",
		"myevent",
		10 * 1000,
	}

	if msg.String() != "id: 123\nretry: 10000\nevent: myevent\ndata: test\n\n" {
		t.Fatal("Messagem not empty.")
	}
}

func TestMultilineDataMessage(t *testing.T) {
	msg := Message{data: "test\ntest"}

	if msg.String() != "data: test\ndata: test\n\n" {
		t.Fail()
	}
}
