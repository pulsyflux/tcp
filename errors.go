package tcp

import "errors"

var (
	errConnectionClosed = errors.New("connection closed")
	errConnectionDead   = errors.New("connection dead")
	errConnectionInUse  = errors.New("connection in use: concurrent Send/Receive not allowed")
	errMessageTooLarge  = errors.New("message exceeds maximum size (16MB)")
)
