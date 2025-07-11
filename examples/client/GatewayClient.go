package main

// import (
// 	"bytes"
// 	"coapws/coap"
// 	"context"
// 	"io"
// 	"log"

// 	"github.com/plgd-dev/go-coap/v3/message/codes"
// 	"github.com/plgd-dev/go-coap/v3/message/pool"
// )

// func main() {
// 	// Connect to the central server
// 	centralClient, err := coap.NewCoAPClient("coap+ws://localhost:8080")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to central server: %v", err)
// 	}
// 	defer centralClient.Close()

// 	// Subscribe to resource updates from central server
// 	observeRequest := pool.NewMessage(context.Background())
// 	observeRequest.SetCode(codes.GET)
// 	observeRequest.SetToken([]byte{0x42})
// 	observeRequest.SetBody(bytes.NewReader([]byte("temperature")))

// 	centralClient.SendMessage(&coap.CoAPMessage{Pool: observeRequest})
// 	log.Println("Subscribed to temperature updates from central server.")
// 	handler := func(msg *coap.CoAPMessage) *coap.CoAPMessage {
// 		log.Printf("Received message: %s", msg.Pool.Body())
// 		response := pool.NewMessage(context.Background())
// 		response.SetCode(codes.Content)
// 		response.SetToken(msg.Pool.Token())
// 		response.SetBody(bytes.NewReader([]byte("Hello from CoAP+WS gateway server!")))
// 		return &coap.CoAPMessage{Pool: response}
// 	}
// 	// Start the gateway server for downstream devices
// 	gatewayServer := coap.NewCoAPServer(handler)

// 	// Register handler for downstream subscriptions
// 	gatewayServer.Router.Handle(0x01, func(msg *coap.CoAPMessage, conn *coap.WebSocketConnection) {
// 		// Read the body from io.ReadSeeker and convert it to a string
// 		body := msg.Pool.Body()
// 		bodyBytes, err := io.ReadAll(body)
// 		if err != nil {
// 			log.Printf("Failed to read message body: %v", err)
// 			return
// 		}
// 		resource := string(bodyBytes) // Assuming payload contains resource identifier
// 		token := msg.Pool.Token()

// 		session, _ := gatewayServer.GetSession(conn)
// 		session.AddObservation(resource, token)

// 		log.Printf("Leaf device subscribed to resource: %s", resource)

// 		// Acknowledge observation
// 		response := pool.NewMessage(context.Background())
// 		response.SetCode(codes.Content)
// 		response.SetToken(token)
// 		response.SetBody(bytes.NewReader([]byte("Observation registered")))
// 		centralClient.SendMessage(&coap.CoAPMessage{Pool: response})

// 	})

// 	// Relay notifications from central server to leaf devices
// 	go func() {
// 		for {
// 			response, err := centralClient.ReceiveMessage()
// 			if err != nil {
// 				log.Fatalf("Error receiving update from central server: %v", err)
// 			}
// 			body := response.Pool.Body()
// 			bodyBytes, err := io.ReadAll(body)
// 			if err != nil {
// 				log.Printf("Failed to read message body: %v", err)
// 				continue
// 			}
// 			gatewayServer.NotifyObservers("temperature", bodyBytes)
// 		}
// 	}()

// 	log.Fatal(gatewayServer.Start(":8081"))
// }
