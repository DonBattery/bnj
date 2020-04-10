package core

import "github.com/gorilla/websocket"

type WS_Conn struct {
	conn *websocket.Conn
}

func NewConn(conn *websocket.Conn) *WS_Conn {
	return &WS_Conn{
		conn: conn,
	}
}
