package util

import (
	"fmt"
	"math"
)

type Position = Vector
type Vector struct {
	X int
	Y int
}

func (p Vector) Add(r Vector) Vector {
	return Vector{p.X + r.X, p.Y + r.Y}
}
func (p Vector) Subtract(r Vector) Vector {
	return Vector{p.X - r.X, p.Y - r.Y}
}
func (p Vector) Dot(r Vector) Vector {
	return Vector{p.X * r.X, p.Y * r.Y}
}
func (p Vector) Scale(c float64) Vector {
	return Vector{int(math.Floor(float64(p.X) * c)), int(math.Floor(float64(p.Y) * c))}
}
func (v Vector) Length() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}
func (v Vector) Equal(r Vector) bool {
	return v.X == r.X && v.Y == r.Y
}
func (v Vector) ToString() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}
