package game

type Rect struct {
	X      float64
	Y      float64
	Width  int
	Height int
}

func NewRect(x, y float64, width, height int) *Rect {
	return &Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (r *Rect) Collide(other *Rect) bool {
	return r.X < other.X+float64(other.Width) &&
		r.X+float64(r.Width) > other.X &&
		r.Y < other.Y+float64(other.Height) &&
		r.Y+float64(r.Height) > other.Y
}

func (r *Rect) CollideMany(others []*Rect) bool {
	for _, other := range others {
		if r.Collide(other) {
			return true
		}
	}
	return false
}
