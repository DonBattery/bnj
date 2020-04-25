package app

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/donbattery/bnj/hub"
	"github.com/donbattery/bnj/server"
)

// Create the run command
func (app *app) runCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"start"},
		Short:   "Spinn up the Bounce 'n Junk server",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(app.ctx)
		},
	}

	runCmd.Flags().IntP("port", "p", 9090, "The PORT where the server will listen")

	return runCmd
}

// run creates and initiates the communication channels, the game, the hub and the http server
func run(ctx context.Context) error {
	// Create the control channel on which the hub will push client control notyfications to the game
	controlCh := make(chan *hub.ControlNotify)
	// Create the game
	game := game.NewGame(ctx, controlCh)
	// Create the hub
	hub := hub.NewHub(ctx, controlCh)
	// Create the server
	server := server.NewServer(ctx)

	// Pass in callback functions the these objects
	hub.SetRequestFn(game.Request)     // the hub can call the game with arbitary client requests (login)
	hub.SetLogoutFn(game.Logout)       // the hub can call the game with when a conn is dropped, to remove the player
	game.SetBroadcastFn(hub.Broadcast) // the game can call the hub to broadcast state update
	server.SetConnectFn(hub.Connect)   // the server can call the hub to add a new WebSocket connection (new client)

	// Start the game and the hub
	game.Start()
	hub.Start()
	// Start the server, return any error
	return server.Start()
}
