package draw

type Pt struct {
	X, Y int
}

func P(x, y int) Pt {
	return Pt{
		X: x,
		Y: y,
	}
}
