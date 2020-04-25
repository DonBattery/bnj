package server

import (
	"errors"

	"github.com/labstack/echo"
)

func (s *Server) hub(c echo.Context) error {
	// Get the Client's ID, return error if not found
	clientId := c.QueryParam("client_id")
	if clientId == "" {
		return errors.New("Empty Client ID")
	}
	// Upgrade the connection to WebSocket, return error if fails
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	// Add the new WebSocket connection (to the Hub)
	s.addConn(clientId, ws)
	return nil
}
