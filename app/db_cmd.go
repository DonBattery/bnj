package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/donbattery/bnj/model"
	"github.com/donbattery/bnj/utils"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			return databaseStats(app.ctx, args)
		},
	}

	databaseCmd.Flags().BoolP("remove", "R", false, "Remove the target, You cannot remove the root or the main buckets")

	return databaseCmd
}

func databaseStats(ctx context.Context, args []string) error {
	if len(args) == 0 {
		args = []string{""}
	}

	db := utils.DB(ctx)

	for _, arg := range args {
		switch db.GetType(arg) {
		case "Bucket":
			// If remove is enabled, and the target is not the root or the main buckets
			if arg != "" && viper.GetBool("remove") {
				if err := db.DeleteBucket(arg); err != nil {
					log.Errorf("Failed to remove bucket %s due to: %s", arg, err.Error())
				} else {
					log.Warnf("Successfully removed bucket %s", arg)
				}
			} else {
				if err := getTree(db, arg); err != nil {
					log.Errorf("Failed to get tree from %s due to %s", arg, err.Error())
				}
			}
		case "Key":
			if viper.GetBool("remove") {
				if err := db.Del(arg); err != nil {
					log.Errorf("Failed to delete key %s due to %s", arg, err.Error())
				} else {
					log.Warnf("Successfully removed key %s", arg)
				}
			} else {
				if err := getKey(db, arg); err != nil {
					log.Errorf("Failed to get key %s due to %s", arg, err.Error())
				}
			}
		default:
			log.Errorf("%s doesn't exists in the Database", arg)
		}
	}

	return nil
}

func getTree(db model.DBConn, keyChain string) error {
	tree, err := db.Tree(keyChain)
	if err != nil {
		return errors.Wrap(err, "Cannot get Database Dump")
	}
	prettyDump, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Cannot marshal Database Dump as JSON")
	}
	fmt.Printf("%s\n", prettyDump)
	return nil
}

func getKey(db model.DBConn, keyChain string) error {
	var value = make(map[string]interface{})

	err := db.Get(keyChain, &value)
	if err != nil {
		return errors.Wrap(err, "Failed to retrieve key")
	}

	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal value %s", value)
	}

	fmt.Printf("%s\n", jsonBytes)

	return nil
}
