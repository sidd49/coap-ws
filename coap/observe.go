package coap

import (
	"context"
	"io"
	"log"
	"strings"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
)

func HandleObserve(msg *CoAPMessage, conn *WebSocketConnection, server *CoAPServer) {
	session, exists := server.GetSession(conn)
	if !exists {
		log.Println("Session not found for connection.")
		return
	}

	// Read the body from io.ReadSeeker and convert it to a string
	body := msg.Pool.Body()
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		log.Printf("Failed to read message body: %v", err)
		return
	}
	resource := string(bodyBytes) // Assuming payload contains resource identifier
	token := msg.Pool.Token()

	if optionExists(msg.Pool.Options(), 6) { // Observe option number is 6
		log.Printf("Adding observation for resource: %s", resource)
		session.AddObservation(resource, token)
		ar := pool.NewMessage(context.Background())
		ar.SetToken(token)
		ar.SetCode(codes.Content)
		ar.SetBody(strings.NewReader("Observation registered"))
		// Send acknowledgment to client
		response := &CoAPMessage{
			Pool: ar,
		}
		conn.SendMessage(response.Encode())
	}
}

func optionExists(options message.Options, optionNumber uint16) bool {
	for _, opt := range options {
		if uint16(opt.ID) == optionNumber {
			return true
		}
	}
	return false
}

func (s *CoAPServer) NotifyObservers(resource string, payload []byte) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	for _, session := range s.Sessions {
		token, exists := session.GetObservationToken(resource)
		if exists {
			ar := pool.NewMessage(context.Background())
			ar.SetToken(token)
			ar.SetCode(codes.Content)
			ar.SetBody(strings.NewReader(string(payload)))
			response := &CoAPMessage{
				Pool: ar,
			}
			session.Conn.SendMessage(response.Encode())
		}
	}
}
