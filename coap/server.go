package coap

import (
	"coapws/coap/models"
	"context"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
)

type CoAPServer struct {
	Upgrader websocket.Upgrader
	Handler  func(*CoAPMessage) *CoAPMessage // Custom message handler
	Router   *Router
	Sessions map[*WebSocketConnection]*Session // Active sessions
	Lock     sync.Mutex
	MongoDb  Mongodb
}

// NewCoAPServer initializes a CoAP+WS server
func NewCoAPServer(handler func(*CoAPMessage) *CoAPMessage) *CoAPServer {
	log.Println("Connecting to MongoDb...")
	mongodb := NewMongo()
	mongodb.Init()
	return &CoAPServer{
		Upgrader: websocket.Upgrader{
			Subprotocols: []string{"coap"},                           // Accept only the "coap" subprotocol
			CheckOrigin:  func(r *http.Request) bool { return true }, // Allow all origins for simplicity
		},
		Handler:  handler,
		Router:   NewRouter(),
		Sessions: make(map[*WebSocketConnection]*Session),
		MongoDb:  *mongodb,
	}
}

// Start starts the CoAP+WS server
func (s *CoAPServer) Start(address string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming WebSocket connection: %s", r.RemoteAddr)

		// Upgrade the HTTP connection to a WebSocket connection
		conn, err := s.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}

		// Check if the selected subprotocol is "coap"
		if conn.Subprotocol() != "coap" {
			log.Printf("Invalid subprotocol: %s", conn.Subprotocol())
			conn.Close()
			return
		}

		log.Println("WebSocket connection established with 'coap' subprotocol.")
		wsConn := &WebSocketConnection{Conn: conn}
		session := s.AddSession(wsConn)
		log.Printf("Session ID: %v", session.ID)

		defer func() {
			s.RemoveSession(wsConn)
			log.Printf("Session %v closed", session.ID)
			wsConn.Close()
		}()
		for {
			data, err := wsConn.ReceiveMessage()
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				break
			}

			msg, err := Decode(data)
			if err != nil {
				log.Printf("Failed to decode message: %v", err)
				continue
			}
			s.Router.Serve(msg, wsConn)

			s.HandlePingPong(msg, wsConn)
			// HandleObserve(msg, wsConn, s)
			log.Printf("Received CoAP message: %v and session ID  %v", msg.Pool.Body(), session.ID)
			deviceid, err := models.ExtractDeviceID(msg.Pool.Body())
			if err != nil || deviceid == "" {
				log.Println("Did not received device id from the iot device, cannot register !!  ", err)
			} else {
				thing := models.Thing{
					SessionID:   strconv.Itoa(int(session.ID)),
					DeviceId:    deviceid,
					Description: "New device",
				}
				var things []models.Thing
				things = append(things, thing)
				// registring / updating data in mongo
				err = s.MongoDb.StoreThings(things)
				if err != nil {
					log.Printf("Could not save the entry in mongo db : %v", err)
					continue
				}
			}

			// Process message using the custom handler
			response := s.Handler(msg)
			if response != nil {
				wsConn.SendMessage(response.Encode())
			}
		}

	})
	log.Printf("CoAP+WS Server started at %s", address)
	return http.ListenAndServe(address, nil)
}

func (s *CoAPServer) HandlePingPong(msg *CoAPMessage, wsConn *WebSocketConnection) {
	if msg.Pool.Code() == codes.Ping { // Ping
		log.Println("Received Ping. Responding with Pong...")

		ar := pool.NewMessage(context.Background())
		hexString, _ := hex.DecodeString("e2922342e6d8c076")
		ar.SetToken(hexString)
		ar.SetCode(codes.Pong)

		pongMessage := &CoAPMessage{
			Pool: ar,
		}

		wsConn.SendMessage(pongMessage.Encode())
	}
}
