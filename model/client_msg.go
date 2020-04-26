package model

import (
	"fmt"
	"strings"
)

type ConnStatus int

const (
	Status_Unknown       ConnStatus = 0
	Status_Connected     ConnStatus = 1
	Status_Authenticated ConnStatus = 2
	Status_InGame        ConnStatus = 3
)

func (cs ConnStatus) StatusIn(statuses []ConnStatus) bool {
	for _, status := range statuses {
		if cs == status {
			return true
		}
	}
	return false
}

func (cs ConnStatus) StatusNotIn(statuses []ConnStatus) bool {
	for _, status := range statuses {
		if cs == status {
			return false
		}
	}
	return true
}

type NotifyType string

const (
	Notify_Chat    NotifyType = "chat"
	Notify_Control NotifyType = "control"
)

// ChatNotify is a client chat notification
type ChatNotify struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
}

// ControlNotify is an user control notification (key down or key up)
type ControlNotify struct {
	ControlType string `json:"control_type"`
	ControlKey  string `json:"control_key"`
}

// ClientNotify is a push message to the hub by a client
type ClientNotify struct {
	NotifyType NotifyType     `json:"notify_type"`
	Chat       *ChatNotify    `json:"chat,omitempty"`
	Control    *ControlNotify `json:"control,omitempty"`
}

// ClientRequest is a request to the hub by a client
type ClientRequest struct {
	ClientId    string `json:"-"`
	RequestId   string `json:"request_id"`
	RequestType string `json:"request_type"`
	RequestBody string `json:"request_body"`
	Response    func(status ServerResponseStatus, payload interface{})
}

// CreateResponse creates the ServerResponse based on the ClientRequest
// With the supplied ServerResponseStatus and payload.
func (cr *ClientRequest) CreateResponse(status ServerResponseStatus, payload string) *ServerResponse {
	return &ServerResponse{
		RequestId:  cr.RequestId,
		Status:     status,
		StatusText: status.String(),
		Payload:    payload,
	}
}

type ClientMsgType string

const (
	ClientMsg_Notify  ClientMsgType = "notify"
	ClientMsg_Request ClientMsgType = "request"
)

// ClientMsg is the message object sent by the client to the server
type ClientMsg struct {
	ClientId      string         `json:"-"`
	ClientMsgType ClientMsgType  `json:"msg_type"`
	Request       *ClientRequest `json:"request,omitempty"`
	Notify        *ClientNotify  `josn:"notify,omitempty"`
}

// ConnError
type ConnError struct {
	connId string
	action string
	code   int
	text   string
}

// NewConnError creates a new ConnError object with the supplied connId, action (which generated the error),
// and the original error-text and/or error
func NewConnError(connId, action string, code int, errs ...interface{}) *ConnError {
	var msg []string
	for _, err := range errs {
		if val, ok := err.(error); ok {
			msg = append(msg, val.Error())
		}
		if val, ok := err.(string); ok {
			msg = append(msg, val)
		}
	}
	return &ConnError{
		connId: connId,
		action: action,
		code:   code,
		text:   strings.Join(msg, ": "),
	}
}

// Error returns the ConnError in the Conn ID: %s Action: %s Error Code: %d Error: %s format
func (ce *ConnError) Error() string {
	return fmt.Sprintf("ConnID: %s Action: %s ErrorCode: %d Error: %s", ce.connId, ce.action, ce.code, ce.text)
}

// ConnId returns the Conn ID
func (ce *ConnError) ConnId() string {
	return ce.connId
}

// Action returns the name of the action which generated the error
func (ce *ConnError) Action() string {
	return ce.action
}

// Code returns the error code
func (ce *ConnError) Code() int {
	return ce.code
}

// Text returns the error-text
func (ce *ConnError) Text() string {
	return ce.text
}
