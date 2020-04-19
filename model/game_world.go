package model

type Game interface {
	Init(serverMsgCh chan *ServerMsg)
	Login(req *ClientRequest) bool
	RemovePlayer(name string)
	// GetWorld() GameWorldDump
}

type GameWorld struct {
	WorldRules   *WorldRules
	Players      []*Player
	WorldMap     *WorldMap
	WorldObjects []*GameObject
}

type GameWorldDump struct {
	WorldRules   WorldRules   `json:"world_rules"`
	Players      []Player     `json:"players"`
	WorldMap     WorldMap     `json:"world_map"`
	WorldObjects []GameObject `json:"world_objects"`
}

func (gw *GameWorld) Dump() GameWorldDump {
	gwd := GameWorldDump{
		WorldRules: *gw.WorldRules,
		WorldMap:   *gw.WorldMap,
	}
	for _, player := range gw.Players {
		gwd.Players = append(gwd.Players, *player)
	}
	for _, obj := range gw.WorldObjects {
		gwd.WorldObjects = append(gwd.WorldObjects, *obj)
	}
	return gwd
}

type WorldRules struct {
	BlockSize   int `josn:"block_size"`
	MaxPlayer   int `json:"max_player"`
	MinPlayer   int `json:"min_player"`
	TargetScore int `json:"target_score"`
	WaitTime    int `json:"wait_time"`
}

type Player struct {
	Name       string `json:"name"`
	Color      string `json:"color"`
	RoundWins  int    `json:"round_wins"`
	RoundScore int    `json:"round_score"`
	TotalScore int    `json:"total_score"`
}

type WorldMap struct {
	Background string   `json:"backgrond"`
	Rows       []string `json:"rows"`
}

type GameObject struct {
	ID       string  `json:"-"`
	ParentID string  `json:"-"`
	Type     string  `json:"type"`
	PosX     float64 `json:"pos_x"`
	PosY     float64 `json:"pos_y"`
	FlipX    bool    `json:"flip_x"`
	FlipY    bool    `json:"flip_y"`
	Width    int     `json:"-"`
	Height   int     `json:"-"`
	Vector   *Vector `json:"-"`
}

type Vector interface {
	X() float64
	Y() float64
	Len() float64
}

func DefaultWorldMap() *WorldMap {
	return &WorldMap{
		Background: "#4d9de3",
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
