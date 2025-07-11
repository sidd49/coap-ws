package coap

const (
	CSM     uint8 = 0x71
	Ping    uint8 = 0x72
	Pong    uint8 = 0x73
	Release uint8 = 0x74
	Abort   uint8 = 0x75
)

// Handle CSM
func HandleCSM(msg *CoAPMessage) {
	// Example: Parse Max-Message-Size or Block-Wise-Transfer capability
}

// Handle Ping/Pong
// func HandlePingPong(msg *CoAPMessage) *CoAPMessage {
//     if uint8(msg.pool.Code()) == Ping {
// 		message := pool.NewMessage(context.Background())
//         return &CoAPMessage{Code: Pong, Token: msg.Token}
//     }
//     return nil
// }
