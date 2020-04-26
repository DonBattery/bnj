package model

type Game interface {
	Init(serverMsgCh chan *ServerMsg)
	Login(req *ClientRequest) bool
	RemovePlayer(name string)
	// GetWorld() GameWorldDump
}

type GameWorldDump struct {
	WorldRules   WorldRules       `json:"world_rules"`
	WorldMap     WorldMap         `json:"world_map"`
	Players      []PlayerDump     `json:"players"`
	WorldObjects []GameObjectDump `json:"world_objects"`
}

type WorldRules struct {
	BlockSize   int     `json:"block_size"   yaml:"block_size"   mapstructure:"block_size"`
	MaxPlayer   int     `json:"max_player"   yaml:"may_player"   mapstructure:"max_player"`
	MinPlayer   int     `json:"min_player"   yaml:"min_player"   mapstructure:"min_player"`
	TargetScore int     `json:"target_score" yaml:"target_score" mapstructure:"target_score"`
	WaitTime    int     `json:"wait_time"    yaml:"wait_time"    mapstructure:"wait_time"`
	Gravity     float64 `json:"gravity"      yaml:"gravity"      mapstructure:"gravity"`
	Friction    float64 `json:"friction"     yaml:"friction"     mapstructure:"friction"`
}

type PlayerDump struct {
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

type GameObjectDump struct {
	ObjType string `json:"obj_type"`
	Anim    int    `json:"anim"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	FlipX   bool   `json:"flip_x"`
	FlipY   bool   `json:"flip_y"`
}

func DefaultWorldMap() WorldMap {
	return WorldMap{
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
