package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/donbattery/bnj/database"
	log "github.com/donbattery/bnj/log"
	"github.com/donbattery/bnj/model"
	"github.com/donbattery/bnj/utils"
)

// app represents the top-level application
type app struct {
	// ctx is the App's context. Used to pass around the configs and the database connection
	ctx context.Context
	// cmd is the root command of the App's CLI
	cmd *cobra.Command
}

// Run instanciates and starts the application
func Run() {
	a := app{}
	a.init()
}

// init sets up and runs the application
func (app *app) init() {
	// Set up the CLI
	app.cmd = app.setupCLI()
	// Silence the default Cobra usage and errors
	app.cmd.SilenceUsage = true
	app.cmd.SilenceErrors = true
	//Execute the CLI
	if err := app.cmd.Execute(); err != nil {
		log.Fatalf("%s", err)
	}
	// Explicitly exit with 0 if no error occured
	os.Exit(0)
}

// setupCLI sets up the root Cobra Command and attaches all subcommands to it
func (app *app) setupCLI() *cobra.Command {
	// The rootCmd command
	rootCmd := &cobra.Command{
		// Get the nice long description
		Long: getLong(),
		// This function runs before every subcommand (before its RunE function is executed)
		// it creates the app's context (with the configs and database connection in it)
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return app.createContext(cmd)
		},
		// This function runs after every subcommand (after its RunE function is executed)
		// it cleans up after the app (closes the db, etc...)
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return app.cleanup()
		},
	}
	// Persistent flags of the CLI, every subcommand will inherit these
	rootCmd.PersistentFlags().String("config_file", "", "Config file")
	rootCmd.PersistentFlags().Bool("debug", false, "Debug Mode")

	// Get all the subcommands and add them to the rootCmd
	rootCmd.AddCommand(
		app.runCmd(),
		// app.versionCmd(),
		// app.initCmd(),
		// app.configCmd(),
		app.databaseCmd(),
		// app.updateCmd(),
		// app.reportCmd(),
	)

	return rootCmd
}

// createContext creates and validates the config object and the database connection
// and injects them into the app's context
func (app *app) createContext(cmd *cobra.Command) error {
	// Get the defaut configurations
	conf := model.DefaultConf()
	// Check flags
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return errors.Wrap(err, "Failed to bind Cobra flags to Viper keys")
	}
	// Check envs
	log.Debug("Searching environment for variables prefixed with BNJ_")
	viper.SetEnvPrefix("BNJ")
	viper.AutomaticEnv()
	// Check conf file
	log.Debug("Checking config file")
	if confPath := viper.GetString("config_file"); confPath != "" {
		viper.SetConfigFile(confPath)
		log.Debugf("config_file is set to %s", confPath)
	} else {
		log.Debug("Searching PWD for the bnj_conf.yaml file")
		viper.SetConfigName("bnj_conf")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrapf(err, "Failed to load configuration file")
	}
	log.Debugf("Configuration file used: %s", viper.ConfigFileUsed())

	// Load the Viper configs into the object
	if err := viper.Unmarshal(&conf); err != nil {
		return errors.Wrap(err, "Failed to decode config object")
	}
	if viper.GetBool("debug") { // If debug is enabled
		niceConf, _ := json.MarshalIndent(conf, "", "  ")
		fmt.Printf("\nConfigs before validation:\n\n%s\n", niceConf) // print the configs
	}
	// Validate the configurations
	if err := conf.Validate(); err != nil {
		return errors.Wrap(err, "Invalid Bounce 'n Junk server configurations")
	}
	log.Debug("Creating context with the valid configurations")
	// Create a new context and put the configs into it
	ctx := context.WithValue(context.Background(), "config", conf)
	// With valid configs we should be able to init the database
	db := database.New()
	if err := db.Init(model.GetDBInitConfig(&conf.DataBase)); err != nil {
		return errors.Wrap(err, "Failed to init the database")
	}
	log.Debugf("Database connection %s@%s injected into the context", conf.DataBase.Type, conf.DataBase.URL)
	// Assign the context to the app, and put the database interface in it as well
	app.ctx = context.WithValue(ctx, "database", db)

	return nil
}

// cleanup closes the database
func (app *app) cleanup() error {
	return utils.DB(app.ctx).Close()
}

// Create the nice long description
func getLong() string {
	return fmt.Sprintf("%s%s",
		color.New(color.FgCyan).Sprint(`
Bounce 'n Junk
`), `
Bounce 'n Junk is a crappy Jump 'n Bump clone.

This self contained application has three major function:
  - Set up the environment for the game server (generate config and db file, build, Dockerize)
  - Run the game server (create game world, serve the game site and handle WebSocket communication)
  - Send admin requests to a running game server (check stats, configs, stop server)
`,
	)
}
