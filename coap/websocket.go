package coap

import (
	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	Conn *websocket.Conn
}

// SendMessage sends a binary message over the WebSocket connection
func (ws *WebSocketConnection) SendMessage(data []byte) error {
	return ws.Conn.WriteMessage(websocket.BinaryMessage, data)
}

// ReceiveMessage reads a binary message from the WebSocket connection
func (ws *WebSocketConnection) ReceiveMessage() ([]byte, error) {
	_, data, err := ws.Conn.ReadMessage()
	return data, err
}

// Close closes the WebSocket connection
func (ws *WebSocketConnection) Close() error {
	return ws.Conn.Close()
}
