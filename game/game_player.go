package game

import "github.com/donbattery/bnj/model"

type player struct {
	clientId   string
	name       string
	color      string
	roundWins  int
	roundScore int
	totalScore int
}

func newPlayer(clientId, name, color string) *player {
	return &player{
		clientId: clientId,
		name:     name,
		color:    color,
	}
}

func (p *player) dump() model.PlayerDump {
	return model.PlayerDump{
		Name:       p.name,
		Color:      p.color,
		RoundWins:  p.roundWins,
		RoundScore: p.roundScore,
		TotalScore: p.totalScore,
	}
}
