package server

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/donbattery/bnj/utils"
)

var (
	upgrader = websocket.Upgrader{}
)

func hub(c echo.Context) error {
	errCounter := 0

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	go func() {
		for {
			if errCounter > 0 {
				errCounter--
			}
			time.Sleep(time.Second * 30)
		}
	}()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Yolo from server $$444!!!1111!!! 1"))
		if err != nil {
			c.Logger().Error(err)
			if errCounter++; errCounter > 6 {
				c.Logger().Error("WebSocket error count is too high, disconnecting...")
				return errors.New("WebSocket error count is too high, disconnecting...")
			}
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			if errCounter++; errCounter > 6 {
				c.Logger().Error("WebSocket error count is too high, disconnecting...")
				return errors.New("WebSocket error count is too high, disconnecting...")
			}
		}
		fmt.Printf("%s\n", msg)
	}
}

func Run(ctx context.Context) error {
	cfg := utils.Conf(ctx)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "frontend")
	e.GET("/hub", hub)

	return e.Start(":" + strconv.Itoa(cfg.Port))
}
