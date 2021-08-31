package game

type Sun struct {
	Orientation int
}

func (s Sun) Move() Sun {
	s.Orientation = (s.Orientation + 1) % 6
	return s
}
