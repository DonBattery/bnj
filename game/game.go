package game

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/donbattery/bnj/hub"
	log "github.com/donbattery/bnj/logger"
	"github.com/donbattery/bnj/utils"
	"github.com/pkg/errors"
)

type Game struct {
	ctx         context.Context
	initOnce    sync.Once
	world       *world
	frame       int64
	step        time.Duration
	controlCh   chan *hub.ControlNotify
	broadcastFn func(msg *hub.ServerMsg)
}

func NewGame(ctx context.Context, step time.Duration, controlCh chan *hub.ControlNotify) *Game {
	cfg := utils.Conf(ctx)

	return &Game{
		ctx:         ctx,
		step:        step,
		frame:       0,
		controlCh:   controlCh,
		world:       createWorld(&cfg.WorldRules, nil),
		broadcastFn: func(msg *ServerMsg) {},
	}
}

func (gc *Game) SetBrouadcastFn(f func(msg *ServerMsg)) {
	gc.broadcastFn = f
}

func (gc *Game) Init() {
	gc.initOnce.Do(func() {
		go gc.run()
	})
}

func (gc *Game) Request(req *hub.ClientRequest) error {
	switch req.RequestType {
	case "login":
		return gc.handleLogin(req)
	}
	resp := fmt.Sprintf("Unknown request type %s", req.RequestType)
	req.Response(hub.ResponseStatusBadRequest, resp)
	return errors.New(resp)
}

func (gc *Game) Logout(clientId string) error {
	return gc.world.removePlayer(clientId)
}

func (gc *Game) run() {
	tick := time.NewTicker(gc.step)
	defer tick.Stop()
	for {
		select { // We need to use double select to avoid the control ch to block the updates
		case <-tick.C:
			gc.frame++
			gc.update()
			gc.broadcastFn(NewServerMsg(hub.ServerMsg_Update, gc.world.objectDump(), nil, nil))
		default:
		}

		select {
		case <-gc.ctx.Done():
			log.Warnf("Game Loop's context is done, returning...")
			return

		case ctl := <-gc.controlCh:
			log.Infof("Incoming control notify %s %s", ctl.ControlKey, ctl.ControlType)

		case <-tick.C:
			gc.frame++
			gc.update()
			gc.broadcastFn(NewServerMsg(hub.ServerMsg_Update, gc.world.objectDump(), nil, nil))
		}
	}
}

func (gc *Game) update() {
	time.Sleep(time.Millisecond * 10)
}

func (gc *Game) broadcast() {
	gc.serverMsgCh <- hub.NewServerMsg(hub.ServerMsg_Update, gc.world.objectDump(), nil, nil)
}

func (gc *Game) handleLogin(req *hub.ClientRequest) error {
	// Validate LoginRequest
	var loginRequest hub.LoginRequest
	if err := json.Unmarshal([]byte(req.RequestBody), &loginRequest); err != nil {
		req.Response(hub.ResponseStatusBadRequest, fmt.Sprintf("Invalid LoginRequest JSON %s", err.Error()))
		return errors.Wrap(err, "Invalid request JSON")
	}

	// Check if a player with the sname name is already connected to the game
	for _, player := range gc.world.players {
		if player.Name == loginRequest.Name {
			resp := fmt.Sprintf("Someone is already connected with the name %s", loginRequest.Name)
			req.Response(hub.ResponseStatusUnauthorized, resp)
			return errors.New(resp)
		}
	}

	// Add the new player
	gc.world.addPlayer(&hub.Player{
		ClientId: req.ClientId,
		Name:     loginRequest.Name,
		Color:    loginRequest.Color,
	})

	// Change the associated wsConn's status to InGame
	utils.Core(gc.ctx).ChangeConnStatus(req.ClientId, hub.Status_InGame)

	// Send the accepted status and the world dump to the player
	req.Response(hub.ResponseStatusAccepted, gc.world.dump())

	return nil
}
