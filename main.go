package main

import (
	"math/rand"
	"time"

	"github.com/donbattery/bnj/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	app.Run()
}
