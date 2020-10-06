package minecraft

import (
	"fmt"
	math2 "github.com/itay2805/mcserver/math"
	"math"
)

type Position struct {
	X, Y, Z	int
}

func (p Position) String() string {
	return fmt.Sprintf("Position{X: %d, Y: %d, Z: %d}", p.X, p.Y, p.Z)
}

func (p Position) ToPoint() math2.Point {
	return math2.NewPoint(float64(p.X), float64(p.Y), float64(p.Z))
}

func ParsePosition(val uint64) Position {
	p := Position{
		X: int((val >> 38) & 0x3FFFFFF),
		Y: int(val & 0xFFF),
		Z: int((val >> 12) & 0x3FFFFFF),
	}

	if p.X > 33554432 {
		p.X -= 67108864
	}

	if p.Z > 33554432 {
		p.Z -= 67108864
	}

	if p.Y > 2048 {
		p.Y -= 4096
	}

	return p
}

func (p Position) Pack() uint64 {
	return uint64(((p.X & 0x3FFFFFF) << 38) | ((p.Z & 0x3FFFFFF) << 12) | (p.Y & 0xFFF))
}

type Angle uint8

func ToAngle(v float32) Angle {
	return Angle(v * 256.0 / 360.0)
}

func (a Angle) ToRadians() float64 {
	return float64(a) * math.Pi / 128.0
}
