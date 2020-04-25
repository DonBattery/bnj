package game

import (
	"testing"

	"github.com/c2fo/testify/require"
)

func Test_Rect(t *testing.T) {
	req := require.New(t)

	t_cases := []struct {
		rect1         *Rect
		rect2         *Rect
		shouldCollide bool
	}{
		{
			rect1:         NewRect(0, 0, 10, 10),
			rect2:         NewRect(0, 0, 1, 1),
			shouldCollide: true,
		},
		{
			rect1:         NewRect(45, 45, 10, 10),
			rect2:         NewRect(54, 54, 1, 1),
			shouldCollide: true,
		},
		{
			rect1:         NewRect(45, 45, 10, 10),
			rect2:         NewRect(55, 54, 10, 10),
			shouldCollide: false,
		},
	}

	for i, t_case := range t_cases {
		t.Logf("Collosion Case %d", i+1)
		req.Equal(t_case.shouldCollide, t_case.rect1.Collide(t_case.rect2), "Collosion should work")
	}
}
