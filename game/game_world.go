package game

import (
	"sync"

	log "github.com/donbattery/bnj/logger"
	"github.com/donbattery/bnj/model"
	"github.com/rs/xid"
)

type gameWorld struct {
	mu       sync.RWMutex
	rules    model.WorldRules
	worldMap model.WorldMap
	players  []*player
	objects  []*gameObject
	rect     *rect
}

func newGameWorld(rules model.WorldRules, worldMap model.WorldMap) *gameWorld {
	return &gameWorld{
		rules:    rules,
		worldMap: worldMap,
		rect:     newRect(0, 0, len(worldMap.Rows[0])*rules.BlockSize, len(worldMap.Rows)*rules.BlockSize),
	}
}

func (gw *gameWorld) dump() model.GameWorldDump {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	var objects []model.GameObjectDump
	for _, obj := range gw.objects {
		objects = append(objects, obj.dump())
	}

	var players []model.PlayerDump
	for _, player := range gw.players {
		players = append(players, player.dump())
	}

	return model.GameWorldDump{
		WorldRules:   gw.rules,
		WorldMap:     gw.worldMap,
		Players:      players,
		WorldObjects: objects,
	}
}

func (gw *gameWorld) objectDump() (objects []model.GameObjectDump) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	for _, obj := range gw.objects {
		objects = append(objects, obj.dump())
	}
	return
}

func (gw *gameWorld) playerDump() (players []model.PlayerDump) {
	gw.mu.RLock()
	defer gw.mu.RUnlock()

	for _, player := range gw.players {
		players = append(players, player.dump())
	}
	return
}

func (gw *gameWorld) addPlayer(p *player) {
	gw.mu.Lock()
	defer gw.spawnChar(p)
	defer gw.mu.Unlock()
	// Add the player to the list of players
	gw.players = append(gw.players, p)
}

func (gw *gameWorld) removePlayer(clientId string) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	found := false

	// Remove the player if he/she is in the game
	for i, player := range gw.players {
		if player.clientId == clientId {
			log.Debugf("Removing player with client ID %s from the game", clientId)
			gw.players = append(gw.players[:i], gw.players[i+1:]...)
			found = true
			break
		}
	}

	// If player is removed remove any belonging child element
	if found {
		for i, obj := range gw.objects {
			if obj.parentId == clientId {
				gw.objects = append(gw.objects[:i], gw.objects[i+1:]...)
			}
		}
	}
}

func (gw *gameWorld) spawnChar(p *player) {
	gw.mu.Lock()
	defer gw.mu.Unlock()

	x, y := gw.findSafePlace(gw.rules.BlockSize)
	gw.objects = append(gw.objects, newGameObject(xid.New().String(), p.clientId, "vita", x, y))
}

func (gw *gameWorld) findSafePlace(size int) (x float64, y float64) {
	x = float64(randInt(0, gw.rect.width-size))
	y = float64(randInt(0, gw.rect.height-size))
	return
}
