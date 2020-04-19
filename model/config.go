package model

import validation "github.com/go-ozzo/ozzo-validation"

// Config is the server config object
type Config struct {
	Port       int        `json:"port"        yaml:"port"        mapstructure:"port"`
	DataBase   DataBase   `json:"database"    yaml:"database"    mapstructure:"database"`
	WorldRules WorldRules `json:"world_rules" yaml:"world_rules" mapstructure:"world_rules"`
}

// Validate the server configurations
func (conf Config) Validate() error {
	return validation.ValidateStruct(&conf,
		validation.Field(&conf.Port, validation.Required, validation.Min(1000), validation.Max(9999)),
	)
}

// DefaultsConf returns the sensible defaut configs
// both the commands and the initial config file uses these
func DefaultConf() Config {
	return Config{
		Port: 9090,
		DataBase: DataBase{
			Type: "bolt",
			URL:  "bnj.db",
		},
		WorldRules: WorldRules{
			BlockSize:   16,
			MaxPlayer:   10,
			MinPlayer:   2,
			TargetScore: 33,
			WaitTime:    90,
		},
	}
}
