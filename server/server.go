package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/donbattery/bnj/model"
	"github.com/donbattery/bnj/utils"
	"github.com/gorilla/websocket"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	ctx       context.Context
	srv       *echo.Echo
	upgrader  websocket.Upgrader
	connectFn func(clientId string, conn *websocket.Conn)
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		ctx:      ctx,
		srv:      echo.New(),
		upgrader: websocket.Upgrader{},
	}
}

// SetConnectFn sets the supplyed function as the server's Connect function
// which will be called with every new Client's ID and WebSocket connection
func (s *Server) SetConnectFn(f func(clientId string, conn *websocket.Conn)) {
	s.connectFn = f
}

// Start sets up and starts the HTTP server
func (s *Server) Start() error {
	// Inject the Configs and the Database into the server's context
	s.srv.Use(newInjectorMiddleware(s.ctx))
	// Use the Logger middleware
	s.srv.Use(middleware.Logger())
	// Use the Recover middleware
	s.srv.Use(middleware.Recover())
	// Serve the Frontend from a static folder TODO: move this out
	s.srv.Static("/", "frontend")
	// Upgrade the requests to /hub route into WebSocket connection
	s.srv.GET("/hub", s.hub)
	// Administrative endpoint
	s.srv.POST("/admin", s.admin)
	// Run the server
	return s.srv.Start(":" + strconv.Itoa(utils.Conf(s.ctx).Port))
}

// newInjectorMiddlewar creates a new Echo middleware that injects the configs
// and the database from the server's context into the Echo context
func newInjectorMiddleware(ctx context.Context) func(echo.HandlerFunc) echo.HandlerFunc {
	cfg := utils.Conf(ctx)
	db := utils.DB(ctx)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", cfg)
			c.Set("database", db)
			return next(c)
		}
	}
}

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
	s.connectFn(clientId, ws)
	return nil
}

func (s *Server) admin(c echo.Context) error {
	cfg := c.Get("config")
	val, ok := cfg.(model.Config)
	if !ok {
		return c.String(http.StatusInternalServerError, "Failed to get config from the context")
	}
	return c.JSON(http.StatusOK, val)
}
