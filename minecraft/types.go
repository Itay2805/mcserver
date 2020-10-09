package minecraft

import (
	"fmt"
	math2 "github.com/itay2805/mcserver/math"
	"log"
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

func (p Position) ApplyFace(face Face) Position {
	switch face {
	case FaceBottom:
		return Position{
			X: p.X,
			Y: p.Y - 1,
			Z: p.Z,
		}
	case FaceTop:
		return Position{
			X: p.X,
			Y: p.Y + 1,
			Z: p.Z,
		}
	case FaceNorth:
		return Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z - 1,
		}
	case FaceSouth:
		return Position{
			X: p.X,
			Y: p.Y,
			Z: p.Z + 1,
		}
	case FaceWest:
		return Position{
			X: p.X - 1,
			Y: p.Y,
			Z: p.Z,
		}
	case FaceEast:
		return Position{
			X: p.X + 1,
			Y: p.Y,
			Z: p.Z,
		}
	}
	log.Panicln("Invalid face", face)
	return Position{}
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

type Face int

const (
	FaceBottom = Face(iota)
	FaceTop
	FaceNorth
	FaceSouth
	FaceWest
	FaceEast
)

func (f Face) Invert() Face {
	switch f {
	case FaceBottom: return FaceTop
	case FaceTop: return FaceBottom
	case FaceNorth: return FaceSouth
	case FaceSouth: return FaceNorth
	case FaceWest: return FaceEast
	case FaceEast: return FaceWest
	default:
		return FaceEast
	}
}

func (f Face) String() string {
	switch f {
	case FaceBottom: return "FaceBottom"
	case FaceTop: return "FaceTop"
	case FaceNorth: return "FaceNorth"
	case FaceSouth: return "FaceSouth"
	case FaceWest: return "FaceWest"
	case FaceEast: return "FaceEast"
	default:
		return fmt.Sprintf("Face(%d)", f)
	}
}


type Shape int

const (
	ShapeStraight = Shape(iota)
	ShapeInnerLeft
	ShapeInnerRight
	ShapeOuterLeft
	ShapeOuterRight
)

func (f Shape) String() string {
	switch f {
	case ShapeStraight: return "ShapeStraight"
	case ShapeInnerLeft: return "ShapeInnerLeft"
	case ShapeInnerRight: return "ShapeInnerRight"
	case ShapeOuterLeft: return "ShapeOuterLeft"
	case ShapeOuterRight: return "ShapeOuterRight"
	default:
		return fmt.Sprintf("Shape(%d)", f)
	}
}


type Hinge int

const (
	HingeLeft = Hinge(iota)
	HingeRight
)

func (f Hinge) String() string {
	switch f {
	case HingeLeft: return "HingeLeft"
	case HingeRight: return "HingeRight"
	default:
		return fmt.Sprintf("Hinge(%d)", f)
	}
}
