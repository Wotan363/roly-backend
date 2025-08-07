package webSocket

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/roly-backend/internal/config"
)

// Defines the HTTP to Websocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return config.AllowedOrigins[origin]
	},
}

// Struct for each Websocket-Connection
type Connection struct {
	ws          *websocket.Conn
	sendChannel chan []byte
	ctx         context.Context
	cancel      context.CancelFunc
	Uuid        string
	cleanupOnce sync.Once // Ensures that cleanup is only executed once, even if multiple goroutines call it concurrently on the same connection.
}

// Handles new incoming websocket connections
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrades initialen HTTP Request to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error while upgrading http connection to websocket connection",
			slog.String("error", err.Error()),
		)
		return
	}

	ctx, cancel := context.WithCancel(context.Background()) // Creates a context so we don't send data to a closed channel when the user disconnected while a message was processed
	// and would have been sent to the client. Because your programm crashes when you try to send data to a closed channel
	conn := &Connection{
		ws:          ws,
		sendChannel: make(chan []byte, 10), // Initialises a buffered channel so multiple simultaniouos messages to the channel won't get lost. Can hold up to 10 messages simultaniously
		ctx:         ctx,
		cancel:      cancel,
		Uuid:        uuid.NewString(),
	}

	slog.LogAttrs(context.Background(), slog.LevelInfo, "New incoming websocket connection",
		slog.String("uuid", conn.Uuid),
	)

	// Starts the Read and Write Loops
	go readLoop(conn)
	go writeLoop(conn)
}

// Reads the messages that the client sends through the websocket connection
func readLoop(conn *Connection) {
	defer cleanup(conn)

	for {
		select {
		case <-conn.ctx.Done():
			return
		default:
			// Read message. If reading takes too long (for example, client istn't available anymore), close connection
			conn.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
			messageType, msg, err := conn.ws.ReadMessage()
			if err != nil {
				// CloseNormalClosure (1000) indicates a clean and intentional shutdown:
				// the client closed the WebSocket connection properly, usually by calling socket.close().

				// CloseGoingAway (1001) means the client is navigating away from the page,
				// closing the tab or window, or otherwise leaving the site,
				// often without explicitly sending a close frame.

				// CloseAbnormalClosure (1006) means the connection was closed abnormally:
				// no close frame was received. This can happen if the browser crashes,
				// the tab is force-closed, or the network connection is suddenly lost.
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.LogAttrs(context.Background(), slog.LevelError, "Error reading incoming websocket message",
						slog.String("uuid", conn.Uuid),
						slog.String("error", err.Error()),
					)
				} else {
					slog.LogAttrs(context.Background(), slog.LevelInfo, "Websocket closed by client",
						slog.String("uuid", conn.Uuid),
					)
				}
				return
			}
			// Handle incoming text messages
			if messageType == websocket.TextMessage {
				if config.DebugMode {
					slog.LogAttrs(context.Background(), slog.LevelDebug, "New incoming websocket message",
						slog.String("uuid", conn.Uuid),
						slog.String("message", string(msg)),
					)
				}
				go handleMessage(conn, msg)
			}
		}
	}
}

// Sends a message to the client through the websocket connection
func writeLoop(conn *Connection) {
	defer cleanup(conn)

	for {
		select {
		case <-conn.ctx.Done():
			return
		case msg, ok := <-conn.sendChannel:
			// if "ok" is true, then the channel has received a message
			// if "ok" is fale, then the channel was closed
			if !ok {
				return
			}
			if config.DebugMode {
				slog.LogAttrs(context.Background(), slog.LevelDebug, "New outgoing websocket message",
					slog.String("uuid", conn.Uuid),
					slog.String("message", string(msg)),
				)
			}
			// Sends the message to the client. If sending takes too long, then close the connection
			conn.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			err := conn.ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				slog.LogAttrs(context.Background(), slog.LevelError, "Error while sending a websocket message to a client",
					slog.String("uuid", conn.Uuid),
					slog.String("error", err.Error()),
					slog.String("message", string(msg)),
				)
				return
			}
		}
	}
}

// Closes and deletes connection correctly after client disconnected
func cleanup(conn *Connection) {
	conn.cleanupOnce.Do(func() {
		conn.cancel()
		conn.ws.Close()
		close(conn.sendChannel)
		slog.LogAttrs(context.Background(), slog.LevelInfo, "Websocket connection was closed",
			slog.String("uuid", conn.Uuid),
		)
	})
}

func StartServer() {
	http.HandleFunc("/ws", handleWebSocket)

	slog.Info(fmt.Sprintf("Server started in %v mode, listening to port %v", config.Env.AppEnv, config.Env.Port))
	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Env.Port), nil)
	slog.Error("Server closed", slog.String("error", err.Error()))
}
