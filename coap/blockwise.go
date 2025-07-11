package coap

type BlockwiseHandler struct {
	MaxBlockSize int
}

func (b *BlockwiseHandler) HandleBlock(msg *CoAPMessage) ([]*CoAPMessage, error) {
	// Split large payloads into blocks using BERT
	msgBodySize, _ := msg.Pool.BodySize()
	if int(msgBodySize) > b.MaxBlockSize {
		// Split logic here
	}
	return []*CoAPMessage{msg}, nil
}
