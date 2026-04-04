package logic

import (
	"math"
	u "simulate/util"
)

type Animal struct {
	Pos   u.Position
	Char  byte
	grid  *[COLUNM][ROW]byte
	trait Trait
}

// act is called each frame
func (a *Animal) act() {
	var dpos u.Vector

	//find nearest predator
	distance := float64(COLUNM*COLUNM + ROW*ROW)
	farthestP := u.Position{X: -1, Y: -1}
	for i := range COLUNM {
		for j := range ROW {
			if a.grid[i][j] == 'P' || a.grid[i][j] == 'p' {
				length := (a.Pos.Subtract(u.Position{X: i, Y: j})).Length()
				if length < distance {
					distance = length
					farthestP = u.Position{X: i, Y: j}
				}
			}
		}
	}

	//found
	if (!farthestP.Equal(u.Position{X: -1, Y: -1}) && distance <= a.trait.fear.val) {
		dpos = moveInDirection(float64(a.Pos.X-farthestP.X), float64(a.Pos.Y-farthestP.Y))
	} else {
		dpos = u.Vector{}
	}

	move(&a.Pos, a.Char, dpos, a.grid)
}
func move(pos *u.Position, char byte, dpos u.Vector, grid *[COLUNM][ROW]byte) {
	newPos := *pos
	newPos = newPos.Add(dpos)

	newPos.X %= COLUNM
	newPos.Y %= ROW
	if newPos.X < 0 {
		newPos.X += COLUNM
	}
	if newPos.Y < 0 {
		newPos.Y += ROW
	}
	if grid[newPos.X][newPos.Y] == '.' {
		grid[pos.X][pos.Y] = '.'
		pos.X = newPos.X
		pos.Y = newPos.Y
		grid[pos.X][pos.Y] = char
	}
}

func moveInDirection(x float64, y float64) u.Vector {
	xSign, ySign := 1, 1
	if x < 0 {
		xSign = -1
	}
	if y < 0 {
		ySign = -1
	}
	//normalize
	length := math.Sqrt(x*x + y*y)
	if length == 0 {
		return u.Vector{}
	}
	x /= length
	y /= length

	xInt := int(math.Ceil(math.Abs(x))) * xSign
	yInt := int(math.Ceil(math.Abs(y))) * ySign
	return u.Vector{X: xInt, Y: yInt}
}
