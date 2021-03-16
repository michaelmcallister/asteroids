package vector

import (
	"fmt"
	"math"
)

type V2D struct {
	X, Y float64
}

func (v V2D) Add(b V2D) V2D {
	return V2D{
		X: v.X + b.X,
		Y: v.Y + b.Y,
	}
}

func (v V2D) Minus(b V2D) V2D {
	return V2D{
		X: v.X - b.X,
		Y: v.Y - b.Y,
	}
}

func (v V2D) Multiply(n float64) V2D {
	return V2D{
		X: v.X * n,
		Y: v.Y * n,
	}
}

func (v V2D) Rotate(degrees float64) V2D {
	// calculate radians
	angle := degrees * math.Pi / 180
	sine := math.Sin(angle)
	cosine := math.Cos(angle)

	// rotation matrix
	matrix := [2][2]float64{{cosine, -sine}, {sine, cosine}}

	return V2D{
		X: matrix[0][0]*v.X + matrix[0][1]*v.Y,
		Y: matrix[1][0]*v.X + matrix[1][1]*v.Y,
	}
}

func (v V2D) Normalize() V2D {
	return v.Divide(v.Magnitude())
}

func (v V2D) Magnitude() float64 {
	c2 := math.Pow(v.X, 2) + math.Pow(v.Y, 2)

	return math.Sqrt(c2)
}

func (v V2D) Divide(n float64) V2D {
	return V2D{
		X: v.X / n,
		Y: v.Y / n,
	}
}

func (v V2D) Limit(n float64) V2D {
	mag := v.Magnitude()
	if mag > n {
		ratio := n / mag
		v.X *= ratio
		v.Y *= ratio
	}
	return v
}

func (v V2D) String() string {
	return fmt.Sprintf("vector{x: %f, y: %f}", v.X, v.Y)
}
