package websocket

import "github.com/kataras/iris/websocket"

type WSClient struct {
	Account string
	ClientType string
	Socket *websocket.Connection
	Send chan []byte
}

