package hub

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	log "github.com/donbattery/bnj/logger"
	"github.com/donbattery/bnj/utils"
)

// Hub is the WebSocket communication controller
type Hub struct {
	// It inherits a context and from that it creates a new one for every conn
	ctx      context.Context
	mu       sync.RWMutex // a Mutex to make the Hub thread safe
	initOnce sync.Once    // Init the Hub only once
	// A list of Connection objects
	conns []*wsConn
	// The channels where the Hub communicates with the connectionsa and the game
	controlCh   chan *bool
	clientMsgCh chan *bool
	errorCh     chan error
	// Two callback functions
	requestFn func(req *ClientRequest) error // one for arbitary client requests (mostly login)
	logoutFn  func(clientId string)          // one for when a connection is closed and the associated player needs to be logged out from the game
}

// NewHub creates a new Hub in the given context and initializes it
func NewHub(ctx context.Context, controlCh chan *bool) *Hub {
	return &Hub{
		ctx:         ctx,
		controlCh:   controlCh,
		clientMsgCh: make(chan *bool),
		errorCh:     make(chan error),
	}
}

func (hub *Hub) SetRequestFn(f func(req *ClientRequest) error) {
	hub.requestFn = f
}

func (hub *Hub) SetLogoutFn(f func(clientId string)) {
	hub.logoutFn = f
}

// Init sets up and starts the WebSocket Hub
func (hub *Hub) Init() {
	hub.initOnce.Do(func() {
		go hub.run()
	})
}

// run constantly selects from 4 possible events: when the context is Done the hub exits,
// when an error is received on a Conn, the conn is closed and the player is removed
// when a client message is received on one of the Conns it gets handled accordingly (in a separate go rutine)
// when a server message is received it is broadcasted to all the clients in the game.
func (hub *Hub) run() {
	for {
		select {
		case <-hub.ctx.Done():
			log.Debug("WebSocket Hub's context is doen, Hub exitting...")
			for _, conn := range hub.conns {
				conn.done()
			}
			return

		case err := <-hub.errorCh:
			if val, ok := err.(*ConnError); ok {
				log.Errorf("Connection error: %s", val.Error())
				hub.removeConn(val.ConnId())
			} else {
				log.Fatalf("Unexpected error occured: %s", err.Error())
			}

		case msg := <-hub.clientMsgCh:
			log.Debugf("Incoming message from Client: %s type: %s", msg.ClientId, msg.ClientMsgType)
			go hub.handleClientMsg(msg)

		case msg := <-hub.serverMsgCh:
			hub.broadcast(msg, Status_InGame)
		}
	}
}

func (hub *Hub) Connect(ws *websocket.Conn, clientId string) {
	connCtx, cancel := context.WithCancel(hub.ctx)
	conn := newWsConn(connCtx, cancel, clientId, ws, hub.clientMsgCh, hub.errorCh)
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.conns = append(hub.conns, conn)
	log.Infof("New client added to the Hub with Id %s", clientId)
}

func (hub *Hub) removeConn(clientId string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for i, conn := range hub.conns {
		if conn.clientId == clientId {
			conn.done()
			hub.conns = append(hub.conns[:i], hub.conns[i+1:]...)
			log.Infof("CLient %s disconnected", conn.clientId)
			if err := utils.Game(hub.ctx).RemovePlayerByClientId(conn.clientId); err != nil {
				log.Errorf("Failed to remove Player by ID %s from the game. Error: %s", conn.clientId, err.Error())
			}
			return
		}
	}
}

func (hub *Hub) ChangeConnStatus(clientId string, status ConnStatus) {
	for _, conn := range hub.conns {
		if conn.clientId == clientId {
			conn.changeStatus(status)
			return
		}
	}
}

func (hub *Hub) notify(clientId string, msg *ServerMsg) {
	log.Debugf("Notifying client %s", clientId)
	for _, conn := range hub.conns {
		if conn.clientId == clientId {
			go conn.snedJson(msg)
			return
		}
	}
}

// createResponder creates a closure around a ClientRequest and returns a function which can respond to that request
func (hub *Hub) createResponder(req *ClientRequest) func(status ServerResponseStatus, payload interface{}) {
	return func(status ServerResponseStatus, payload interface{}) {
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
				status = ResponseStatusServerError
			}
		}
		// Create the ServerResponse
		resp := req.CreateResponse(status, string(payloadBytes))
		// Create the ServerMsg and put the response in it
		msg := &ServerMsg{
			MsgType:     ServerMsg_Response,
			GameObjects: nil,
			Chat:        nil,
			Response:    resp,
		}
		// Send the response to the requesting client
		log.Debugf("Responding to client: %s request: %s status: %s", req.ClientId, req.RequestId, status.String())
		hub.notify(req.ClientId, msg)
	}
}

func (hub *Hub) broadcast(msg *ServerMsg, statuses ...ConnStatus) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to encode ServerMessage as JSON %s", err.Error())
	}
	preparedMsg, err := websocket.NewPreparedMessage(1, jsonBytes)
	if err != nil {
		log.Fatalf("Failed to create Prepared Websocket Message %s", err.Error())
	}
	for _, conn := range hub.conns {
		go func(c *wsConn, statuses ...ConnStatus) {
			if len(statuses) > 0 {
				if c.status.StatusNotIn(statuses) {
					return
				}
			}
			c.sendPrepared(preparedMsg)
		}(conn, statuses...)
	}
}

func (hub *Hub) handleClientMsg(msg *ClientMsg) {
	switch msg.ClientMsgType {
	// in case of Notification
	case ClientMsg_Notify:
		switch msg.Notify.NotifyType {
		// in case of Chat
		case Notify_Chat:
			hub.onChat(msg.Notify.Chat)
		// in case of Control
		case Notify_Control:
			msg.Notify.Control.ClientId = msg.ClientId // copy the ClientId into the ControlNotify
			hub.controlCh <- msg.Notify.Control
		}

	// in case of Request
	case ClientMsg_Request:
		msg.Request.ClientId = msg.ClientId                     // copy the ClientId into the response
		msg.Request.Response = hub.createResponder(msg.Request) // create responder and put it into the request
		if err := utils.Game(hub.ctx).Request(msg.Request); err != nil {
			log.Errorf("Failed to process request ID: %s Error: %s", msg.Request.ClientId, err.Error())
		}
	}
}

func (hub *Hub) onChat(chat *ChatNotify) {
	log.Debugf("channel: %s msg: %s", chat.Channel, chat.Message)
	msg := NewServerMsg(ServerMsg_Chat, nil, chat, nil)
	hub.broadcast(msg, Status_Authenticated, Status_InGame)
}
