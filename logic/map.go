package logic

import (
	"fmt"
	"math/rand"
	"simulate/util"
)

const ROW = 16
const COLUNM = 8

type Creature interface {
	act()
}

var grid [COLUNM][ROW]byte
var creatures []Creature

func AddAnimals(num int) {
	for range num {
		//get random ascii char, reserve p for predator
		var char byte
		for {
			char = byte(rand.Intn(94) + 33)
			if char != '.' && char != 'P' && char != 'p' {
				break
			}
		}

		//get random coordinate
		var x, y int
		for {
			x = rand.Intn(COLUNM)
			y = rand.Intn(ROW)
			if grid[x][y] == '.' {
				break
			}
		}
		//Add
		AddAnimal(Animal{Pos: util.Position{X: x, Y: y}, Char: char})
	}
}
func AddAnimal(a Animal) {
	a.grid = &grid
	a.trait.vision.init()
	a.trait.fear.init()
	creatures = append(creatures, &a)
	move(&a.Pos, a.Char, util.Vector{}, a.grid)
}
func AddPredator(p Predator) {
	p.grid = &grid
	creatures = append(creatures, &p)
	move(&p.Pos, 'P', util.Vector{}, p.grid)
}
func PrintMap() {
	for i := range creatures {
		creatures[i].act()
	}

	for i := range COLUNM {
		for j := range ROW {
			fmt.Printf("%c", grid[i][j])
		}
		fmt.Println()
	}

}

func InitGrid() {
	for i := range COLUNM {
		for j := range ROW {
			grid[i][j] = '.'
		}
	}
}
