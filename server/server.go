package server

import (
	"context"
	"strconv"

	"github.com/donbattery/bnj/game"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/donbattery/bnj/core"
	"github.com/donbattery/bnj/server/middlewares"
	"github.com/donbattery/bnj/server/routes"
	"github.com/donbattery/bnj/utils"
)

// Run sets up and sruns the HTTP server in the given context
func Run(ctx context.Context) error {
	// Create the GameWrold
	bnjGame := game.NewGameController(ctx)
	// Inject the bjnGame into a new context
	coreCtx := context.WithValue(ctx, "game", bnjGame)
	// Create the WebSocket Core HUB in that context and Init it
	core := core.NewWsHub(coreCtx)
	core.Init()
	// Create the HTTP Server
	server := echo.New()
	// Inject the Core into the server's context
	server.Use(middlewares.NewInjectorMiddleware(coreCtx, core))
	// Use HTTP logger middleware
	server.Use(middleware.Logger())
	// Use the Recover midleware
	server.Use(middleware.Recover())
	// Serve the Frontend from a static folder TODO: move this out
	server.Static("/", "frontend")
	// Upgrade the requests to /hub route into WebSocket connection
	server.GET("/hub", routes.Hub)
	// Administrative endpoint
	server.POST("/admin", routes.Admin)
	// Any other request will be routed here
	// server.Any("*", routes.NotFound)
	// Run the server
	return server.Start(":" + strconv.Itoa(utils.Conf(ctx).Port))
}
