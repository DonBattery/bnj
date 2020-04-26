package utils

import (
	"context"

	"github.com/labstack/echo"

	log "github.com/donbattery/bnj/logger"
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

// EchoConf gets the config object from an Echo context
func EchoConf(e echo.Context) model.Config {
	conf, ok := e.Get("config").(model.Config)
	if !ok {
		log.Fatal("Failed to get the configs from the Echo context")
	}
	return conf
}

// EchoDB gets the database from an Echo context
func EchoDB(e echo.Context) model.DBConn {
	db, ok := e.Get("database").(model.DBConn)
	if !ok {
		log.Fatal("Failed to get the database from the Echo context")
	}
	return db
}
