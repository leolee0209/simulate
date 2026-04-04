package logic

import (
	"math/rand"
)

type Trait struct {
	vision Vision
	fear   Fear
}

// the range which is able to see predator
type Vision struct {
	max float64
	min float64
	val float64
}

// the range which starts to run from predator
type Fear struct {
	max float64
	min float64
	val float64
}

func (v *Vision) init() {
	v.max = min(COLUNM, ROW) / 1.5
	v.min = min(COLUNM, ROW) / 6
	v.val = rand.Float64()*(v.max-v.min) + v.min
}
func (v *Fear) init() {
	v.max = min(COLUNM, ROW) / 1.5
	v.min = min(COLUNM, ROW) / 6
	v.val = rand.Float64()*(v.max-v.min) + v.min
}
