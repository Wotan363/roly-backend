package webSocket

import "fmt"

// Handles the incoming messages and what to do with them (basically like an api endpoint)
func handleMessage(conn *Connection, msg []byte) {

	// Just as an example a echo
	select {
	case <-conn.ctx.Done():
		return
	case conn.sendChannel <- fmt.Appendf(nil, "%s", msg):
		// everything is fine
	}
}
