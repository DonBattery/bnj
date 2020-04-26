package game

type rect struct {
	x      float64
	y      float64
	width  int
	height int
}

func newRect(x, y float64, width, height int) *rect {
	return &rect{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

func (r *rect) collide(other *rect) bool {
	return r.x < other.x+float64(other.width) &&
		r.x+float64(r.width) > other.x &&
		r.y < other.y+float64(other.height) &&
		r.y+float64(r.height) > other.y
}

func (r *rect) collideMany(others []*rect) bool {
	for _, other := range others {
		if r.collide(other) {
			return true
		}
	}
	return false
}
