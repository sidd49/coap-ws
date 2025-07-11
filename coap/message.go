package coap

import (
	"context"
	"fmt"

	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/tcp/coder"
)

type CoAPMessage struct {
	Pool *pool.Message
}

// Encode encodes the CoAP message for WebSocket transport
func (msg *CoAPMessage) Encode() []byte {
	response, err := msg.Pool.MarshalWithEncoder(coder.DefaultCoder)
	if err != nil {
		return nil
	}
	return response
}

// Decode decodes a WebSocket binary frame into a CoAPMessage
func Decode(data []byte) (*CoAPMessage, error) {
	msg := pool.NewMessage(context.Background())
	_, err := msg.UnmarshalWithDecoder(coder.DefaultCoder, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CoAP message: %v", err)
	}
	return &CoAPMessage{Pool: msg}, nil
}
