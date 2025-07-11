package main

// import (
// 	"bufio"
// 	"coapws/coap"
// 	"context"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/plgd-dev/go-coap/v3/message"
// 	"github.com/plgd-dev/go-coap/v3/message/codes"
// 	"github.com/plgd-dev/go-coap/v3/message/pool"
// )

// func main() {
// 	coapWSURL := "coap+ws://localhost:8080"

// 	client, err := coap.NewCoAPClient(coapWSURL)
// 	if err != nil {
// 		log.Fatalf("Failed to create CoAP client: %v", err)
// 	}
// 	defer client.Close()

// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Println("Enter messages to send to the server. Type 'exit' to quit.")

// 	for {
// 		fmt.Print("Enter message: ")
// 		userInput, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Fatalf("Failed to read input: %v", err)
// 		}
// 		userInput = strings.TrimSpace(userInput)
// 		if userInput == "exit" {
// 			break
// 		}

// 		ar := pool.NewMessage(context.Background())
// 		hexString, _ := hex.DecodeString("e2922342e6d8c076")
// 		ar.SetToken(hexString)
// 		ar.SetCode(codes.GET)
// 		ar.SetContentFormat(message.TextPlain)
// 		ar.SetBody(strings.NewReader(userInput))
// 		// Create a CoAP message
// 		msg := &coap.CoAPMessage{Pool: ar}

// 		// Send the message
// 		if err := client.SendMessage(msg); err != nil {
// 			log.Printf("Failed to send message: %v", err)
// 			continue
// 		}
// 		log.Printf("Sent message: %s", msg.Pool.Body())

// 		// Receive a response
// 		response, err := client.ReceiveMessage()
// 		if err != nil {
// 			log.Printf("Failed to receive message: %v", err)
// 			continue
// 		}
// 		body, err := io.ReadAll(response.Pool.Body())
// 		if err != nil {
// 			log.Printf("Failed to read body: %v", err)
// 			continue
// 		}
// 		log.Printf("Received response: %s", string(body))
// 	}

// 	log.Println("Client has exited.")
// }
