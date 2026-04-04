package main

import (
	"fmt"
	"simulate/cmd"
	"simulate/logic"
	"simulate/util"
	"time"
)

func main() {
	logic.InitGrid()

	logic.AddPredator(logic.Predator{Pos: util.Position{X: 5, Y: 5}})
	logic.AddAnimals(3)

	for i := 0; ; i++ {
		fmt.Println()
		fmt.Println()

		logic.PrintMap()

		time.Sleep(1 * time.Second)
		cmd.CallClear()
	}
}
