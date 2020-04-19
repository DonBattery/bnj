package game

import (
	"context"
	"sync"

	"github.com/donbattery/bnj/model"
	"github.com/donbattery/bnj/utils"
)

type GameController struct {
	ctx         context.Context
	mu          sync.RWMutex
	initOnce    sync.Once
	world       *model.GameWorld
	serverMsgCh chan *model.ServerMsg
}

func NewGameController(ctx context.Context) *GameController {
	cfg := utils.Conf(ctx)

	return &GameController{
		ctx: ctx,
		world: &model.GameWorld{
			WorldRules: &cfg.WorldRules,
			WorldMap:   model.DefaultWorldMap(),
		},
	}
}

func (gc *GameController) Init(serverMsgCh chan *model.ServerMsg) {
	gc.initOnce.Do(func() {
		gc.serverMsgCh = serverMsgCh
	})
}

func (gc *GameController) GetWorld() model.GameWorld {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	return *gc.world
}

func (gc *GameController) SetMap(worldMap *model.WorldMap) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.world.WorldMap = worldMap
}

func (gc *GameController) AddPlayer(player *model.Player) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.world.Players = append(gc.world.Players, player)
}

func (gc *GameController) RemovePlayer(name string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	for i, player := range gc.world.Players {
		if player.Name == name {
			gc.world.Players = append(gc.world.Players[:i], gc.world.Players[i+1:]...)
			return
		}
	}
}

func (gc *GameController) AddObject(objects ...*model.GameObject) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	gc.world.WorldObjects = append(gc.world.WorldObjects, objects...)
}

func (gc *GameController) RemoveObject(ids ...string) {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	var newObjects []*model.GameObject
	for _, object := range gc.world.WorldObjects {
		if utils.Excludes(ids, object.ID) {
			newObjects = append(newObjects, object)
		}
	}

	gc.world.WorldObjects = newObjects
}
