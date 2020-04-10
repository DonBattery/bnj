package utils

import (
	"context"

	log "github.com/donbattery/bnj/log"
	"github.com/donbattery/bnj/model"
)

// Conf gets the config object from context
func Conf(ctx context.Context) model.Config {
	conf, ok := ctx.Value("config").(model.Config)
	if !ok {
		log.Fatal("Failed to get the configs from the context")
	}
	return conf
}

// DB gets the database from the context
func DB(ctx context.Context) model.DBConn {
	db, ok := ctx.Value("database").(model.DBConn)
	if !ok {
		log.Fatal("Failed to get the database from the context")
	}
	return db
}
