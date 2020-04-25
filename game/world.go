package game

import (
	"math"
	"sync"

	log "github.com/donbattery/bnj/logger"
	"github.com/donbattery/bnj/model"
	"github.com/pkg/errors"
)

type world struct {
	mu       sync.RWMutex
	rules    *model.WorldRules
	players  []*model.Player
	worldMap *model.WorldMap
	elements []*GameElement
	rect     *Rect
}

func createWorld(rules *model.WorldRules, worldMap *model.WorldMap) *world {
	if worldMap == nil {
		worldMap = model.DefaultWorldMap()
	}
	return &world{
		rules:    rules,
		worldMap: worldMap,
		rect:     NewRect(0, 0, worldMap.Width()*rules.BlockSize, worldMap.Height()*rules.BlockSize),
	}
}

func (w *world) addPlayer(player *model.Player) {
	w.mu.Lock()
	// defer w.spawnCharacter(player)
	defer w.mu.Unlock()
	// Add the player to the list of players
	w.players = append(w.players, player)
}

func (w *world) spwanCharacter(player *model.Player) {
	w.mu.Lock()
	defer w.mu.Unlock()

	character := NewCharacter(player)

	x, y := w.findSafePlace(character.Size())

	character.PosX = x
	character.PosY = y

	// Create and add the player's Character to the world elements
	w.elements = append(w.elements, character)
}

func (w *world) findSafePlace(width, height int) (posX, posY float64) {
	return 0, 0
}

// removePlayer removes a player from the game by ClientId
// and all associated child element (player controlled character)
func (w *world) removePlayer(clientId string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	found := false

	// Remove the player if he/she is in the game
	for i, player := range w.players {
		if player.ClientId == clientId {
			log.Debugf("Removing player with client ID %s from the game", clientId)
			w.players = append(w.players[:i], w.players[i+1:]...)
			found = true
			break
		}
	}

	// If player is removed remove any belonging child element
	if found {
		for i, elem := range w.elements {
			if elem.Parent() == clientId {
				w.elements = append(w.elements[:i], w.elements[i+1:]...)
			}
		}
		return nil
	}

	// If not found return error
	return errors.Errorf("Cannot remove player. There is no player in the game with clientId %s", clientId)
}

func (w *world) dump() model.GameWorldDump {
	return model.GameWorldDump{
		WorldRules:  *w.rules,
		WorldMap:    *w.worldMap,
		Players:     w.playerDump(),
		GameObjects: w.objectDump(),
	}
}

func (w *world) objectDump() (out []model.GameObject) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, elem := range w.elements {
		out = append(out, model.GameObject{
			ObjType: elem.ObjType,
			Anim:    elem.Anim,
			Effect:  elem.Effect,
			PosX:    math.Round(elem.PosX),
			PosY:    math.Round(elem.PosY),
			FlipX:   elem.FlipX,
			FlipY:   elem.FlipY,
		})
	}
	return
}

func (w *world) playerDump() (out []model.Player) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, player := range w.players {
		out = append(out, *player)
	}
	return
}
