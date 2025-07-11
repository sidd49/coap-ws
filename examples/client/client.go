package main

import (
	"bytes"
	"coapws/coap"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
)

type RequestBody struct {
	DeviceID string `json:"deviceid"`
}

func main() {
	coapWSURL := "coap+ws://localhost:8080"
	numClients := 10 // Number of clients to spawn
	var wg sync.WaitGroup
	wg.Add(numClients)

	for i := 0; i < numClients; i++ {
		if i%1000 == 0 {
			time.Sleep(1000 * time.Millisecond)
		}

		go func(clientID int) {
			defer wg.Done()

			client, err := coap.NewCoAPClient(coapWSURL)
			if err != nil {
				log.Printf("Client %d: Failed to create CoAP client: %v", clientID, err)
				return
			}
			defer client.Close()

			ar := pool.NewMessage(context.Background())
			hexString, _ := hex.DecodeString("e2922342e6d8c076")
			ar.SetToken(hexString)
			ar.SetCode(codes.GET)
			ar.SetContentFormat(message.AppJSON)
			id := uuid.New()
			stringID := id.String()
			var reqBody RequestBody = RequestBody{
				DeviceID: stringID,
			}
			jsonData, err := json.Marshal(reqBody)
			if err != nil {
				log.Fatalf("Error marshalling JSON: %v", err)
			}
			log.Printf("Marshal data : %v", string(jsonData))
			finalBody := bytes.NewReader(jsonData)
			log.Printf("Final Body : %v", finalBody)
			ar.SetBody(finalBody)
			// Create a CoAP message
			msg := &coap.CoAPMessage{Pool: ar}

			// Send the message
			if err := client.SendMessage(msg); err != nil {
				log.Printf("Client %d: Failed to send message: %v", clientID, err)
				return
			}
			log.Printf("Client %d: Sent message: %s", clientID, msg.Pool.Body())

			time.Sleep(5 * time.Second)
			// Send a PING message
			client.SendPing([]byte{0x01, 0x02, 0x03, 0x04})

			// Receive a response
			response, err := client.ReceiveMessage()
			if err != nil {
				log.Printf("Client %d: Failed to receive message: %v", clientID, err)
				return
			}
			body, err := io.ReadAll(response.Pool.Body())
			if err != nil {
				log.Printf("Client %d: Failed to read body: %v", clientID, err)
				return
			}
			log.Printf("Client %d: Received response: %s", clientID, string(body))
		}(i)
	}

	wg.Wait()
	log.Println("All clients have finished.")
}
