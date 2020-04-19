package game

import (
	"encoding/json"
	"fmt"

	"github.com/donbattery/bnj/model"
)

func (gc *GameController) Login(req *model.ClientRequest) bool {
	var loginRequest model.LoginRequest
	if err := json.Unmarshal([]byte(req.RequestBody), &loginRequest); err != nil {
		req.Response(model.ResponseStatusBadRequest, fmt.Sprintf("Invalid LoginRequest JSON %s", err.Error()))
		return false
	}
	for _, player := range gc.world.Players {
		if player.Name == loginRequest.Name {
			req.Response(model.ResponseStatusUnauthorized, fmt.Sprintf("Someone is already connected with the name %s", loginRequest.Name))
			return false
		}
	}
	gc.AddPlayer(&model.Player{
		Name:  loginRequest.Name,
		Color: loginRequest.Color,
	})

	req.Response(model.ResponseStatusAccepted, gc.world.Dump())

	return true
}
