package game

type vector struct {
	x float64
	y float64
}

func newVector(x, y float64) *vector {
	return &vector{
		x: x,
		y: y,
	}
}

func (vec *vector) X(values ...float64) float64 {
	if len(values) > 0 {
		vec.x = values[0]
	}
	return vec.x
}

func (vec *vector) Y(values ...float64) float64 {
	if len(values) > 0 {
		vec.y = values[0]
	}
	return vec.y
}
