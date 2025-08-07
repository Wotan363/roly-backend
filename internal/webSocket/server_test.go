package webSocket

import (
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestWebsocket(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	header := http.Header{}
	header.Set("Origin", "http://localhost:3000")

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		t.Fatalf("Error while connecting to websocket: %v", err)
	}
	defer c.Close()

	// Sends a message
	err = c.WriteMessage(websocket.TextMessage, []byte("ping"))
	if err != nil {
		t.Fatalf("Error sending message to websocket: %v", err)
	}

	// Reads the answer
	_, message, err := c.ReadMessage()
	if err != nil {
		t.Fatalf("Error reading websocket answer: %v", err)
	}

	log.Printf("Answer from the server: %s", message)

	// Checks if answer was as expected
	if string(message) != "ping" {
		t.Errorf("Unexpected answer from server: %s", message)
	}
}
