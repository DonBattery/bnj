package app

import (
	"github.com/spf13/cobra"

	"github.com/donbattery/bnj/app/controller"
)

// Create the database Command
func (app *app) databaseCmd() *cobra.Command {
	databaseCmd := &cobra.Command{
		Use:     "db",
		Aliases: []string{"database"},
		Short:   "Bounce 'n Junk server database info",
		Long: `
Discover and manipulate the Bounce 'n Junk server database

You can specify multiple chains of database keys (joined with ".") as arguments for this subcommand
(e.g.: $ bnj db users.joe)
If no argument defined, the root bucket will be checked.

If the chain points to a bucket, a tree will be displayd from that bucket
to every sub-bucket and keys, but only the size of the keys will be displayed, not their value.
If the chain points to a key its value will be displyed.

With the --remove flag you may delete buckets or keys, but not the root or the main-buckets

The database subcommand displays valid JSON which can be further parsed (for example with jq)
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return controller.DatabaseStats(app.ctx, args)
		},
	}

	databaseCmd.Flags().BoolP("remove", "R", false, "Remove the target, You cannot remove the root or the main buckets")

	return databaseCmd
}
