package logic

import (
	"math/rand"
)

type Trait struct {
	vision            Vision
	roaming           Roaming
	selfishHerdChance SelfishHerdChance
}

// the range which is able to see predator
type Vision struct {
	max float64
	min float64
	val float64
}

// the movement weight used while roaming toward visible prey
type Roaming struct {
	max float64
	min float64
	val float64
}

// chance to blend escape direction toward nearby prey during evasion
type SelfishHerdChance struct {
	max float64
	min float64
	val float64
}

func (v *Vision) init() {
	v.max = 20
	v.min = 5
	v.val = rand.Float64()*(v.max-v.min) + v.min
}

func (r *Roaming) init() {
	r.max = 1
	r.min = 0
	r.val = rand.Float64()*(r.max-r.min) + r.min
}

func (s *SelfishHerdChance) init() {
	s.max = 1
	s.min = 0
	s.val = clamp01(rand.Float64()/2+0.5)
}
