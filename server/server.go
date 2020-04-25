package server

import (
	"context"
	"strconv"

	"github.com/donbattery/bnj/utils"
	"github.com/gorilla/websocket"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type Server struct {
	ctx      context.Context
	srv      *echo.Echo
	upgrader websocket.Upgrader
	addConn  func(clientId string, conn *websocket.Conn)
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		ctx:      ctx,
		srv:      echo.New(),
		upgrader: websocket.Upgrader{},
	}
}

// SetAddConnFn sets the supplyed function as the server's addConn function
// which will be called with every new Client's ID and WebSocket connection
func (s *Server) SetAddConnFn(f func(clientId string, conn *websocket.Conn)) {
	s.addConn = f
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
	s.srv.Static("/", "_frontend")
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
