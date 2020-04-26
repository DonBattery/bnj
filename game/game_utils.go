package game

import "math/rand"

func randInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
