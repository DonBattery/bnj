package hub

import (
	"fmt"

	"github.com/donbattery/bnj/game"
)

// ServerMsgType is the type of the ServerMsg, sent to a client
// From this field the client can conclude which field of the ServerMsg needs to be processed
type ServerMsgType string

const (
	ServerMsg_Chat     ServerMsgType = "chat"
	ServerMsg_Response ServerMsgType = "response"
	ServerMsg_Update   ServerMsgType = "update"
)

// ServerResponseStatus is sent to a client along with a response to a request
// it is similar to HTTP statuses
type ServerResponseStatus int

// The named ServerResponseStatus
const (
	ResponseStatusOK           ServerResponseStatus = 200
	ResponseStatusAccepted     ServerResponseStatus = 202
	ResponseStatusBadRequest   ServerResponseStatus = 400
	ResponseStatusUnauthorized ServerResponseStatus = 401
	ResponseStatusServerError  ServerResponseStatus = 500
)

// String translates a ServerResponseStatus number to its string representation
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
	MsgType     ServerMsgType   `json:"msg_type"`
	GameObjects []game.Object   `josn:"game_objects,omitempty"`
	Chat        *ChatNotify     `json:"chat,omitempty"`
	Response    *ServerResponse `json:"response,omitempty"`
}

// NewServerMsg creates a ServerMsg with the given type, game-objects, chat-notify and response
// The type determines which field will be checked by the client
func NewServerMsg(msgType ServerMsgType, gameObjects []GameObject, chat *ChatNotify, response *ServerResponse) *ServerMsg {
	return &ServerMsg{
		MsgType:     msgType,
		GameObjects: gameObjects,
		Chat:        chat,
		Response:    response,
	}
}
