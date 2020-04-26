package core

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	log "github.com/donbattery/bnj/logger"
	"github.com/donbattery/bnj/model"
)

// WsHub is the WebSocket communication controller
type WsHub struct {
	ctx      context.Context
	mu       sync.RWMutex
	initOnce sync.Once

	conns []*wsConn

	clientMsgCh chan *model.ClientMsg
	errorCh     chan error
	controlCh   chan *model.ControlNotify

	requestFn func(req *model.ClientRequest)
	logoutFn  func(clientId string)
}

// NewWsHub creates a new WsHub in the given context and initializes it
func NewWsHub(ctx context.Context, controlCh chan *model.ControlNotify) *WsHub {
	return &WsHub{
		ctx:         ctx,
		clientMsgCh: make(chan *model.ClientMsg),
		errorCh:     make(chan error),
		controlCh:   controlCh,
		requestFn:   func(req *model.ClientRequest) {},
	}
}

func (hub *WsHub) SetRequestFn(f func(req *model.ClientRequest)) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.requestFn = f
}

func (hub *WsHub) SetLogoutFn(f func(clientId string)) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.logoutFn = f
}

// Start sets up and starts the WebSocket Hub
func (hub *WsHub) Start() {
	hub.initOnce.Do(func() {
		go hub.reader()
	})
}

// Reading incoming messages on the wsHub's pushQueue and dispatching it to all Conns
func (hub *WsHub) reader() {
	for {
		select {
		case <-hub.ctx.Done():
			log.Debug("WebSocket Hub's context is doen, WsHub exitting...")
			for _, conn := range hub.conns {
				conn.done()
			}
			return

		case err := <-hub.errorCh:
			if val, ok := err.(*model.ConnError); ok {
				log.Errorf("Connection error: %s", val.Error())
				hub.removeConn(val.ConnId())
			} else {
				log.Fatalf("Unexpected error occured: %s", err.Error())
			}

		case msg := <-hub.clientMsgCh:
			log.Debugf("Incoming message from Client: %s type: %s", msg.ClientId, msg.ClientMsgType)
			go hub.handleClientMsg(msg)
		}

	}
}

func (hub *WsHub) Connect(clientId string, ws *websocket.Conn) {
	connCtx, cancel := context.WithCancel(hub.ctx)
	conn := newWsConn(connCtx, cancel, clientId, ws, hub.clientMsgCh, hub.errorCh)
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.conns = append(hub.conns, conn)
	log.Infof("New client added to the Hub with ID %s", clientId)
}

func (hub *WsHub) removeConn(clientId string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for i, conn := range hub.conns {
		if conn.clientId == clientId {
			conn.done()
			hub.conns = append(hub.conns[:i], hub.conns[i+1:]...)
			log.Infof("CLient %s disconnected", conn.clientId)
			go hub.logoutFn(clientId)
		}
	}
}

func (hub *WsHub) ChangeConnStatus(clientId string, status model.ConnStatus) {
	for _, conn := range hub.conns {
		if conn.clientId == clientId {
			conn.changeStatus(status)
			return
		}
	}
}

func (hub *WsHub) notify(clientId string, msg *model.ServerMsg) {
	log.Debugf("Notifying client %s", clientId)
	for _, conn := range hub.conns {
		if conn.clientId == clientId {
			go conn.snedJson(msg)
			return
		}
	}
}

// createResponder creates a closure around a ClientRequest and returns a function which can respond to that request
func (hub *WsHub) createResponder(req *model.ClientRequest) func(status model.ServerResponseStatus, payload interface{}) {
	return func(status model.ServerResponseStatus, payload interface{}) {
		var (
			payloadBytes []byte
			err          error
		)
		if val, ok := payload.(string); ok { // if payload is already a string no need to marshal
			payloadBytes = []byte(val)
		} else { // Marshal the payload as JSON
			payloadBytes, err = json.Marshal(payload)
			if err != nil { // on marshal error, change the status and the payload
				log.Errorf("Failed marshal ServerResponse Payload as JSON: %s", err.Error())
				payloadBytes = []byte("JSON Marshal error")
				status = model.ResponseStatusServerError
			}
		}
		// Create the ServerResponse
		resp := req.CreateResponse(status, string(payloadBytes))
		// Create the ServerMsg and put the response in it
		msg := &model.ServerMsg{
			MsgType:     model.ServerMsg_Response,
			WorldUpdate: nil,
			Chat:        nil,
			Response:    resp,
		}
		// Send the response to the requesting client
		log.Debugf("Responding to client: %s request: %s status: %s", req.ClientId, req.RequestId, status.String())
		hub.notify(req.ClientId, msg)
	}
}

func (hub *WsHub) Broadcast(msg *model.ServerMsg, statuses ...model.ConnStatus) {
	// log.Debugf("Broadcasting message type %s to status: %+v", msg.Type, statuses)
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to encode ServerMessage as JSON %s", err.Error())
	}
	preparedMsg, err := websocket.NewPreparedMessage(1, jsonBytes)
	if err != nil {
		log.Fatalf("Failed to create Prepared Websocket Message %s", err.Error())
	}
	for _, conn := range hub.conns {
		go func(c *wsConn, statuses ...model.ConnStatus) {
			if len(statuses) > 0 {
				if c.status.StatusNotIn(statuses) {
					return
				}
			}
			c.sendPrepared(preparedMsg)
		}(conn, statuses...)
	}
}

func (hub *WsHub) BroadcastGameUpdate(msg *model.ServerMsg) {
	hub.Broadcast(msg, model.Status_InGame)
}

func (hub *WsHub) handleClientMsg(msg *model.ClientMsg) {
	switch msg.ClientMsgType {
	// in case of Notification
	case model.ClientMsg_Notify:
		switch msg.Notify.NotifyType {
		// in case of Chat
		case model.Notify_Chat:
			hub.onChat(msg.Notify.Chat)
		// in case of Control
		case model.Notify_Control:
			hub.onControl(msg.Notify.Control)
		}

	// in case of Request
	case model.ClientMsg_Request:
		msg.Request.ClientId = msg.ClientId                     // copy the client id into the response
		msg.Request.Response = hub.createResponder(msg.Request) // create responder and put it into the request
		go hub.requestFn(msg.Request)
	}
}

func (hub *WsHub) onChat(chat *model.ChatNotify) {
	log.Debugf("channel: %s msg: %s", chat.Channel, chat.Message)
	msg := model.NewServerMsg(model.ServerMsg_Chat, nil, chat, nil)
	hub.Broadcast(msg, model.Status_Authenticated, model.Status_InGame)
	// hub.Broadcast(msg, model.Status_Connected, model.Status_Authenticated, model.Status_InGame)
}

func (hub *WsHub) onControl(control *model.ControlNotify) {

}
