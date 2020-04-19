package core

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"

	log "github.com/donbattery/bnj/log"
	"github.com/donbattery/bnj/model"
)

// wsConn represents a WebSocket connection with a client
type wsConn struct {
	// The context of the connection
	ctx context.Context
	// Close the connection
	done context.CancelFunc
	// Embed Gorilla's WebSocket implementation
	*websocket.Conn
	// A mutex to protect the connection from concurrent writes
	mu sync.Mutex
	// The unique ID of the WebSocket client
	clientId string
	// The status of the connection
	status model.ConnStatus
	// Init the wsConn only once
	initOnce sync.Once
	// Messages from the client will be pushed to this channel
	msgCh chan *model.ClientMsg
	// wsConn errors will be pushed to this channel
	errorCh chan error
}

// newWsConn creates a new wsConn object in the given context, with the given Client ID, WebSocket connection
// and with a supplyed client message and an error channel. Then inits the conn and returns it.
func newWsConn(ctx context.Context, done context.CancelFunc, clientId string, ws *websocket.Conn, msgCh chan *model.ClientMsg, errorCh chan error) *wsConn {
	// create
	conn := &wsConn{
		ctx:      ctx,
		done:     done,
		Conn:     ws,
		clientId: clientId,
		status:   model.Status_Connected,
		msgCh:    msgCh,
		errorCh:  errorCh,
	}
	// init
	conn.initOnce.Do(func() {
		go conn.clientMsgReader()
	})
	// return
	return conn
}

// changeStatus changes the wsConn's status to the supplied ConnStatus
func (conn *wsConn) changeStatus(status model.ConnStatus) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	conn.status = status
}

// sendRaw sends the supplied byte slice to the client as a TextMessage
func (conn *wsConn) sendRaw(msg []byte) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	err := conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		connErr := model.NewConnError(conn.clientId, "Write Raw", -1, err)
		conn.errorCh <- connErr
	}
}

// snedJosn sends the supplied message to the client as JSON
func (conn *wsConn) snedJson(msg interface{}) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	err := conn.WriteJSON(msg)
	if err != nil {
		connErr := model.NewConnError(conn.clientId, "Write JSON", -1, err)
		conn.errorCh <- connErr
	}
}

// sendPrepared sends the supplied prepared websocket message to the client
func (conn *wsConn) sendPrepared(msg *websocket.PreparedMessage) {
	conn.mu.Lock()
	defer conn.mu.Unlock()
	err := conn.WritePreparedMessage(msg)
	if err != nil {
		connErr := model.NewConnError(conn.clientId, "Write Prepared Message", -1, err)
		conn.errorCh <- connErr
	}
}

// clientMsgReader constantly tries to read the next incoming message on the wsConn
// upon read it determines the type of the WebSocket message and calls the
// appropriate handler. when the wsConn's context is done, reader returns releasing the resoureces
func (conn *wsConn) clientMsgReader() {
	log.Debugf("Setting up WebSocket reader for ClientID %s", conn.clientId)

	for {
		select {
		case <-conn.ctx.Done():
			log.Debugf("WebSocket %s context is done, releasing connection...", conn.clientId)
			return
		default:
			msgType, msgData, msgErr := conn.ReadMessage()
			if msgErr != nil {
				connErr := model.NewConnError(conn.clientId, "Read Message", -2, msgErr)
				conn.errorCh <- connErr
			}
			log.Debugf("Incoming WebSocket message from client: %s type: %d", conn.clientId, msgType)
			switch msgType {
			case -1:
				// on error do nothing, as we handled it before
			case 1:
				conn.processClientMsg(msgData)
			case 8:
				log.Warnf("Client %s is cloesing the connection %s", conn.clientId, msgData)
			default:
				log.Warnf("Unknown WebSocket message type %d Message: %s", msgType, msgData)
			}
		}
	}
}

// processClientMsg unmarshals a Client Message, received on the wsConn, and pushes it to the
// wsConn's msgCh channel. If an unmarshal error occures it will be pushed to the errorCh channel
func (conn *wsConn) processClientMsg(msgData []byte) {
	msg := &model.ClientMsg{
		ClientId: conn.clientId,
	}
	if err := json.Unmarshal(msgData, msg); err != nil {
		connErr := model.NewConnError(conn.clientId, "Unmarshal", -3, err)
		conn.errorCh <- connErr
		return
	}
	conn.msgCh <- msg
}
