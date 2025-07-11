package coap

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
)

type CoAPClient struct {
	WSConnection *WebSocketConnection
}

// NewCoAPClient initializes a CoAP+WS client
func NewCoAPClient(coapWSURL string) (*CoAPClient, error) {
	log.Printf("Connecting to CoAP+WS server: %s", coapWSURL)

	parsedURL, err := url.Parse(coapWSURL)
	if err != nil {
		log.Fatalf("Failed to parse URL: %v", err)
		return nil, err
	}

	// Translate "coap+ws" to "ws" and "coaps+ws" to "wss"
	switch parsedURL.Scheme {
	case "coap+ws":
		parsedURL.Scheme = "ws"
	case "coaps+ws":
		parsedURL.Scheme = "wss"
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
	}

	// Perform WebSocket handshake with "coap" subprotocol
	conn, _, err := websocket.DefaultDialer.Dial(parsedURL.String(), map[string][]string{
		"Sec-WebSocket-Protocol": {"coap"},
	})
	if err != nil {
		log.Printf("Failed to establish WebSocket connection: %v", err)
		return nil, err
	}

	log.Printf("WebSocket connection established with subprotocol: %s", conn.Subprotocol())
	if conn.Subprotocol() != "coap" {
		log.Fatalf("Server did not accept the 'coap' subprotocol. Got: %s", conn.Subprotocol())
		conn.Close()
		return nil, fmt.Errorf("Invalid subprotocol")
	}

	return &CoAPClient{WSConnection: &WebSocketConnection{Conn: conn}}, nil
}

// SendMessage sends a CoAP message
func (c *CoAPClient) SendMessage(msg *CoAPMessage) error {
	return c.WSConnection.SendMessage(msg.Encode())
}

// ReceiveMessage waits for a CoAP response
func (c *CoAPClient) ReceiveMessage() (*CoAPMessage, error) {
	data, err := c.WSConnection.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	return Decode(data)
}

func (c *CoAPClient) SendPing(token []byte) error {
	ar := pool.NewMessage(context.Background())
	hexString, _ := hex.DecodeString("e2922342e6d8c076")
	ar.SetToken(hexString)
	ar.SetCode(codes.Ping)
	pingMessage := &CoAPMessage{
		Pool: ar,
	}

	log.Println("Sending Ping...")
	return c.WSConnection.SendMessage(pingMessage.Encode())
}

// Close closes the WebSocket connection
func (c *CoAPClient) Close() error {
	return c.WSConnection.Close()
}
