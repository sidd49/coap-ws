package main

import (
	"bytes"
	"coapws/coap"
	"context"
	"log"
	"time"

	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
)

func main() {
    // Define a simple handler function
    handler := func(msg *coap.CoAPMessage) *coap.CoAPMessage {
        log.Printf("Received message: %s", msg.Pool.Body())
        response := pool.NewMessage(context.Background())
        response.SetCode(codes.Content)
        response.SetToken(msg.Pool.Token())
        response.SetBody(bytes.NewReader([]byte("Hello from CoAP+WS server!")))
        return &coap.CoAPMessage{Pool: response}
    }

    // Initialize the CoAP+WS server with the handler
    server := coap.NewCoAPServer(handler)

    // Register observe handler
    server.Router.Handle(0x01, func(msg *coap.CoAPMessage, conn *coap.WebSocketConnection) { 
        coap.HandleObserve(msg, conn, server)
    })

    // Notify devices about resource updates
    go func() {
        for {
            time.Sleep(5 * time.Second)
            server.NotifyObservers("temperature", []byte("Temperature: 30°C"))
        }
    }()

    log.Fatal(server.Start(":8080"))
}