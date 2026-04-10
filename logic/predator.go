package logic

import (
	u "simulate/util"
)

type Predator struct {
	Pos           u.Position
	vision        Vision
	mode          BehaviorMode
	baseColor     ColorRGB
	restTicksLeft int
	dead          bool
}

const predatorSpeed = preySpeed * 3

const (
	predatorRestTicks   = 120
	predatorCatchRadius = 0.5
)

func (p *Predator) act(w *World) {
	if p.restTicksLeft > 0 {
		p.mode = RestingMode
		p.restTicksLeft--
		return
	}

	nearestPrey, found, _ := w.nearestPrey(p.Pos)
	if !found  {
		p.mode = RestingMode
		p.restTicksLeft = predatorRestTicks
		return
	}

	p.mode = ChasingMode
	toPrey := w.wrappedVector(p.Pos, nearestPrey)
	dpos := moveWithSpeed(toPrey, predatorSpeed)
	w.move(&p.Pos, 'P', dpos)

	if w.eatPreyAt(p.Pos, predatorCatchRadius) {
		p.mode = RestingMode
		p.restTicksLeft = predatorRestTicks
		return
	}
}

func (p *Predator) snapshot() CreatureSnapshot {
	return CreatureSnapshot{Pos: p.Pos, Char: 'P', Kind: PredatorKind, Mode: p.mode, BaseColor: p.baseColor}
}

func (p *Predator) alive() bool {
	return !p.dead
}
