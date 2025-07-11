package coap

import "log"

// Router is responsible for routing CoAP messages to appropriate handlers
type Router struct {
    handlers map[uint8]func(*CoAPMessage, *WebSocketConnection) // Handlers for each message code
}

// NewRouter creates and initializes a new Router
func NewRouter() *Router {
    return &Router{
        handlers: make(map[uint8]func(*CoAPMessage, *WebSocketConnection)),
    }
}

// Handle registers a handler for a specific CoAP code
func (r *Router) Handle(code uint8, handler func(*CoAPMessage, *WebSocketConnection)) {
    r.handlers[code] = handler
}

// Serve routes the CoAP message to the appropriate handler
func (r *Router) Serve(msg *CoAPMessage, conn *WebSocketConnection) {
    if handler, exists := r.handlers[uint8(msg.Pool.Code())]; exists {
        handler(msg, conn) // Call the registered handler
    } else {
        log.Printf("No handler found for code: %d", msg.Pool.Code())
    }
}
