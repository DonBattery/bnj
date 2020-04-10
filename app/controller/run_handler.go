package controller

import (
	"context"

	"github.com/donbattery/bnj/server"
)

func Run(ctx context.Context) error {
	return server.Run(ctx)
}
