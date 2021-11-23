package main

type (
	// Compass - an enum with enumerants N=0, NE, E, SW, S, SW, W, NW, CENTER
	Compass uint8

	// Coord - signed int16 pair, absolute location or difference of locations
	Coord struct {
		X, Y int
	}

	Polar struct {
		mag       int
		Direction Compass
	}
)

const (
	N Compass = iota
	NE
	E
	SE
	S
	SW
	W
	NW
	Center
)

func (c Coord) isNormalized() bool {
	return c.X >= -1 && c.X <= 1 && c.Y >= -1 && c.Y <= 1
}

// Rotate a Compass value by the specified number of steps. There are
// eight steps per full rotation. Positive values are clockwise; negative
// values are counterclockwise. E.g., rotate(4) returns a direction 90
// degrees to the right.
func (c Compass) Rotate(n int) Compass {
	return Compass((int(c) + n) % int(Center))
}

var compassNames = map[Compass]string{
	S:      "S",
	SE:     "SE",
	E:      "E",
	NE:     "NE",
	N:      "N",
	NW:     "NW",
	W:      "W",
	SW:     "SW",
	Center: "Center",
}

func (c Compass) String() string {
	return compassNames[c]
}

func (c Compass) asNormalizedCoord() Coord {
	switch c {
	case N:
		return Coord{X: 0, Y: 1}
	case NE:
		return Coord{X: 1, Y: 1}
	case E:
		return Coord{X: 1, Y: 0}
	case SE:
		return Coord{X: 1, Y: -1}
	case S:
		return Coord{X: 0, Y: -1}
	case SW:
		return Coord{X: -1, Y: -1}
	case W:
		return Coord{X: -1, Y: 0}
	case NW:
		return Coord{X: -1, Y: 1}
	default:
		panic(42)
	}
}
