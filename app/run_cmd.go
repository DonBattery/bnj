package app

import (
	"github.com/spf13/cobra"

	"github.com/donbattery/bnj/app/controller"
)

// Create the run command
func (app *app) runCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:     "run",
		Aliases: []string{"start"},
		Short:   "Spinn up the Bounce 'n Junk server",
		Long:    ``,
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.Run(app.ctx)
		},
	}

	runCmd.Flags().IntP("port", "p", 9090, "The PORT where the server will listen")

	return runCmd
}
