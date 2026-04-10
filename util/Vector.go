package util

import (
	"fmt"
	"math"
)

type Position = Vector
type Vector struct {
	X float64
	Y float64
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
	return Vector{p.X * c, p.Y * c}
}
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func (v Vector) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}
func (v Vector) Equal(r Vector) bool {
	return v.X == r.X && v.Y == r.Y
}
func (v Vector) ToString() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}
