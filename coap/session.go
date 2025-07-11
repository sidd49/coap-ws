package coap

import (
	"sync"
	"sync/atomic"
)

var SessionID int32 = 0

type Session struct {
	ID           int32                // Unique identifier for the session
	Conn         *WebSocketConnection // WebSocket connection for the session
	Observations map[string][]byte    // Map of observed resources and associated tokens
	Lock         sync.Mutex           // Ensure thread-safe access
}

// AddObservation adds a new observation to the session
func (s *Session) AddObservation(resource string, token []byte) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Observations[resource] = token
}

// RemoveObservation removes an observation from the session
func (s *Session) RemoveObservation(resource string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Observations, resource)
}

// GetObservationToken retrieves the token for an observed resource
func (s *Session) GetObservationToken(resource string) ([]byte, bool) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	token, exists := s.Observations[resource]
	return token, exists
}

func (s *CoAPServer) AddSession(conn *WebSocketConnection) *Session {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	session := &Session{
		ID:           atomic.AddInt32(&SessionID, 1),
		Conn:         conn,
		Observations: make(map[string][]byte),
	}
	s.Sessions[conn] = session
	return session
}

// RemoveSession removes a session
func (s *CoAPServer) RemoveSession(conn *WebSocketConnection) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	delete(s.Sessions, conn)
}

// GetSession retrieves a session by connection
func (s *CoAPServer) GetSession(conn *WebSocketConnection) (*Session, bool) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	session, exists := s.Sessions[conn]
	return session, exists
}
