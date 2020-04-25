package game

import (
	"github.com/donbattery/bnj/utils"
	"github.com/rs/xid"
)

type WorldRules struct {
	BlockSize   int     `json:"block_size"`
	MaxPlayer   int     `json:"max_player"`
	MinPlayer   int     `json:"min_player"`
	TargetScore int     `json:"target_score"`
	WaitTime    int     `json:"wait_time"`
	Gravity     float64 `json:"gravity"`
	Friction    float64 `json:""friction`
}

type Player struct {
	ClientId   string `json:"-"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	RoundWins  int    `json:"round_wins"`
	RoundScore int    `json:"round_score"`
	TotalScore int    `json:"total_score"`
}

type WorldMap struct {
	Background string   `json:"background"`
	Rows       []string `json:"rows"`
}

func (wm WorldMap) Width() int {
	if wm.Rows != nil {
		return len(wm.Rows[0])
	}
	return 0
}

func (wm WorldMap) Height() int {
	return len(wm.Rows)
}

type GameObject struct {
	ObjType string  `json:"obj_type"`
	Anim    int     `json:"anim"`
	Effect  string  `json:"effect"`
	PosX    float64 `json:"pos_x"`
	PosY    float64 `json:"pos_y"`
	FlipX   bool    `json:"flip_x"`
	FlipY   bool    `json:"flip_y"`
}

type GameWorldDump struct {
	WorldRules  WorldRules   `json:"world_rules"`
	Players     []Player     `json:"players"`
	WorldMap    WorldMap     `json:"world_map"`
	GameObjects []GameObject `json:"game_objects"`
}

func DefaultWorldMap() *WorldMap {
	return &WorldMap{
		Background: "#1e3b69",
		Rows: []string{
			"1110000000000000000000",
			"1000000000001000011000",
			"1000111100001100000000",
			"1000000000011110000011",
			"1100000000111000000001",
			"1110001111110000000001",
			"1000000000000011110001",
			"1000000000000000000011",
			"1110011100000000000111",
			"1000000000003100000001",
			"1000000000031110000001",
			"1011110000311111111001",
			"1000000000000000000001",
			"1100000000000000000011",
			"2222222214000001333111",
			"1111111111111111111111",
		},
	}
}

type GameElement struct {
	obj      *GameObject
	elemId   string
	parentId string
	width    int
	height   int
	vector   *Vector
}

func NewElement(id, parentId, objType, effect string, anim int, posX, posY float64, flipX, flipY bool) *GameElement {
	width, height := utils.SizeFromObjType(objType)
	return &GameElement{
		obj: &GameObject{
			ObjType: objType,
			Effect:  effect,
			Anim:    anim,
			PosX:    posX,
			PosY:    posY,
			FlipX:   flipX,
			FlipY:   flipY,
		},
		elemId:   id,
		parentId: parentId,
		width:    width,
		height:   height,
		vector:   NewVector(0, 0),
	}
}

func NewCharacter(player *Player) *GameElement {
	return NewElement(xid.New().String(), player.ClientId, "bunny", "", 0, 0, 0, false, false)
}

func (ge *GameElement) ID() string {
	return ge.elemId
}

func (ge *GameElement) Parent() string {
	return ge.parentId
}

func (ge *GameElement) Size() (width, height int) {
	return ge.width, ge.height
}
