package logic

import (
	u "simulate/util"
)

type Predator struct {
	Pos  u.Position
	grid *[COLUNM][ROW]byte
}

func (p *Predator) act() {
	var dpos u.Vector

	//find nearest predator
	distance := float64(COLUNM*COLUNM + ROW*ROW)
	farthestP := u.Position{X: -1, Y: -1}
	for i := range COLUNM {
		for j := range ROW {
			if p.grid[i][j] != '.' && p.grid[i][j] != 'P' {
				length := (p.Pos.Subtract(u.Position{X: i, Y: j})).Length()
				if length < distance {
					distance = length
					farthestP = u.Position{X: i, Y: j}
				}
			}
		}
	}

	//found
	if (!farthestP.Equal(u.Position{X: -1, Y: -1})) {
		dpos = moveInDirection(-float64(p.Pos.X-farthestP.X), -float64(p.Pos.Y-farthestP.Y))
		println("found prey", farthestP.ToString())
		println("persuit", dpos.ToString())
	} else {
		dpos = u.Vector{}
	}

	move(&p.Pos, 'P', dpos,p.grid)
}
