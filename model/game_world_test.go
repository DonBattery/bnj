package model

import (
	"testing"

	"github.com/c2fo/testify/require"
)

func Test_GetFloat(t *testing.T) {
	req := require.New(t)

	worldMap := DefaultWorldMap()

	t_cases := []struct {
		x        float64
		y        float64
		size     int
		required int
	}{
		{
			x:        0,
			y:        0,
			size:     0,
			required: 49,
		},
		{
			x:        999,
			y:        0,
			size:     16,
			required: 49,
		},
		{
			x:        0,
			y:        9,
			size:     16,
			required: 49,
		},
		{
			x:        50,
			y:        0,
			size:     16,
			required: 48,
		},
		{
			x:        0,
			y:        0,
			size:     16,
			required: 49,
		},
		{
			x:        4,
			y:        14,
			size:     1,
			required: 50,
		},
		{
			x:        37,
			y:        27,
			size:     3,
			required: 51,
		},
		{
			x:        -1,
			y:        0,
			size:     3,
			required: 49,
		},
		{
			x:        600,
			y:        0,
			size:     3,
			required: 49,
		},
		{
			x:        0,
			y:        -900,
			size:     3,
			required: 48,
		},
	}

	for _, t_case := range t_cases {
		req.Equal(
			t_case.required,
			worldMap.GetFloat(t_case.x, t_case.y, t_case.size),
			"GetFloat should return %d when the block size is %d X is %f and Y is %f", t_case.required, t_case.size, t_case.x, t_case.y)
	}
}
