package draw

type Vec struct {
	P1, P2 Pt
}

func V(p1, p2 Pt) Vec {
	return Vec{
		P1: p1,
		P2: p2,
	}
}
