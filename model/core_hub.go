package model

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Core is the interface of the WebSocket controller
type Core interface {
	AddConn(ws *websocket.Conn, clientID string)
	Broadcast(msg *ServerMsg, statuses ...ConnStatus)
}

type ServerMsgType string

const (
	ServerMsg_Chat     ServerMsgType = "chat"
	ServerMsg_Response ServerMsgType = "response"
	ServerMsg_Update   ServerMsgType = "update"
)

type ServerResponseStatus int

const (
	ResponseStatusOK           ServerResponseStatus = 200
	ResponseStatusAccepted     ServerResponseStatus = 202
	ResponseStatusBadRequest   ServerResponseStatus = 400
	ResponseStatusUnauthorized ServerResponseStatus = 401
	ResponseStatusServerError  ServerResponseStatus = 500
)

func (srs *ServerResponseStatus) String() string {
	switch *srs {
	case ResponseStatusOK:
		return "Response Status: OK"
	case ResponseStatusAccepted:
		return "Response Status: Accepted"
	case ResponseStatusBadRequest:
		return "Response Status: Bad Request"
	case ResponseStatusUnauthorized:
		return "Response Status: Unauthorized"
	case ResponseStatusServerError:
		return "Response Status: Server Error"
	default:
		return fmt.Sprintf("Response Status: unknown status: %d", srs)
	}
}

// ServerResponse is a respons to a ClientRequest
type ServerResponse struct {
	RequestId  string               `json:"request_id"`
	Status     ServerResponseStatus `json:"status"`
	StatusText string               `json:"status_text"`
	Payload    string               `json:"payload"`
}

// ServerMsg is an object to be sent to one or more clients
type ServerMsg struct {
	Type     ServerMsgType   `json:"type"`
	Objects  []GameObject    `josn:"objects,omitempty"`
	Chat     *ChatNotify     `json:"chat,omitempty"`
	Response *ServerResponse `json:"response,omitempty"`
}

func NewServerMsg(msgType ServerMsgType, objects []GameObject, chat *ChatNotify, response *ServerResponse) *ServerMsg {
	return &ServerMsg{
		Type:     msgType,
		Objects:  objects,
		Chat:     chat,
		Response: response,
	}
}
