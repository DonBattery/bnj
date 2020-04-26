package game

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/donbattery/bnj/logger"
	"github.com/donbattery/bnj/model"
	"github.com/donbattery/bnj/utils"
)

type GameController struct {
	mu           sync.RWMutex
	ctx          context.Context
	initOnce     sync.Once
	world        *gameWorld
	frame        int64
	step         time.Duration
	controlCh    chan *model.ControlNotify
	broadcastFn  func(msg *model.ServerMsg)
	connStatusFn func(clientId string, status model.ConnStatus)
}

func NewGameController(ctx context.Context, step time.Duration, controlCh chan *model.ControlNotify) *GameController {
	cfg := utils.Conf(ctx)

	return &GameController{
		ctx:          ctx,
		step:         step,
		frame:        0,
		controlCh:    controlCh,
		world:        newGameWorld(cfg.WorldRules, model.DefaultWorldMap()),
		broadcastFn:  func(msg *model.ServerMsg) {},
		connStatusFn: func(clientId string, status model.ConnStatus) {},
	}
}

//////////////////////
/// Public Methods //
////////////////////

func (gc *GameController) SetBroadcastFn(f func(msg *model.ServerMsg)) {
	gc.broadcastFn = f
}

func (gc *GameController) SetConnStatusFn(f func(clientId string, status model.ConnStatus)) {
	gc.connStatusFn = f
}

func (gc *GameController) Start() {
	gc.initOnce.Do(func() {
		go gc.run()
	})
}

func (gc *GameController) Request(req *model.ClientRequest) {
	switch req.RequestType {
	case "login":
		gc.handleLogin(req)
	default:
		req.Response(model.ResponseStatusBadRequest, fmt.Sprintf("Unknown request type %s", req.RequestType))
	}
}

func (gc *GameController) Logout(clientId string) {
	gc.world.removePlayer(clientId)
}

///////////////////////
/// Private Methods //
/////////////////////

func (gc *GameController) run() {
	tick := time.NewTicker(gc.step)
	defer tick.Stop()

	for {
		select { // We need to use double select to avoid the control ch to block the updates
		case <-tick.C:
			gc.frame++
			gc.update()
			go gc.broadcastFn(model.NewServerMsg(
				model.ServerMsg_Update,
				&model.WorldUpdate{
					Players:      gc.world.playerDump(),
					WorldObjects: gc.world.objectDump(),
				}, nil, nil))
		default: // With an empty default we proceed to the next select without stucking here
		}

		select {
		case <-gc.ctx.Done():
			log.Warnf("Game Loop's context is done, returning...")
			return

		case ctl := <-gc.controlCh:
			log.Infof("Incoming control notify %s %s", ctl.ControlKey, ctl.ControlType)

		case <-tick.C: // once again we check for the ticker
			gc.frame++
			gc.update()
			go gc.broadcastFn(model.NewServerMsg(
				model.ServerMsg_Update,
				&model.WorldUpdate{
					Players:      gc.world.playerDump(),
					WorldObjects: gc.world.objectDump(),
				}, nil, nil))
		}
	}
}

func (gc *GameController) update() {
	time.Sleep(time.Millisecond * 10)
}

func (gc *GameController) handleLogin(req *model.ClientRequest) {
	// Validate LoginRequest
	var loginRequest model.LoginRequest
	if err := json.Unmarshal([]byte(req.RequestBody), &loginRequest); err != nil {
		req.Response(model.ResponseStatusBadRequest, fmt.Sprintf("Invalid LoginRequest JSON %s", err.Error()))
		return
	}

	if len(gc.world.players) >= gc.world.rules.MaxPlayer {
		req.Response(model.ResponseStatusNotAccaptable, "Server is full")
		return
	}

	// Check if a player with the sname name is already connected to the game
	for _, player := range gc.world.players {
		if player.name == loginRequest.Name {
			resp := fmt.Sprintf("Someone is already connected with the name %s", loginRequest.Name)
			req.Response(model.ResponseStatusUnauthorized, resp)
			return
		}
	}

	// Add the new player
	gc.world.addPlayer(newPlayer(req.ClientId, loginRequest.Name, loginRequest.Color))

	// Change the associated wsConn's status to InGame
	gc.connStatusFn(req.ClientId, model.Status_InGame)

	// Send the accepted status and the world dump to the player
	req.Response(model.ResponseStatusAccepted, gc.world.dump())
}
