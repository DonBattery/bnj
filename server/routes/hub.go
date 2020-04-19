package routes

import (
	"errors"

	"github.com/donbattery/bnj/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrader = websocket.Upgrader{}
)

func Hub(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	clientID := c.QueryParam("client_id")
	if clientID == "" {
		return errors.New("Empty Client ID")
	}
	utils.EchoCore(c).AddConn(ws, clientID)
	return nil
}
