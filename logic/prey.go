package logic

import (
	"math"
	"math/rand"
	u "simulate/util"
)

type Prey struct {
	Pos       u.Position
	Char      byte
	trait     Trait
	heading   u.Vector
	mode      BehaviorMode
	baseColor ColorRGB
	dead      bool
}

const preySpeed = 0.15
const preyTurnJitterDegrees = 20.0

// act is called each frame
func (a *Prey) act(w *World) {
	if a.heading.Length() == 0 {
		a.heading = randomUnitVector()
	}
	a.heading = jitterDirection(a.heading, preyTurnJitterDegrees)

	nearestP, found, distance := w.nearestPredator(a.Pos)

	baseDirection := a.heading
	if found && distance <= a.trait.vision.val {
		a.mode = EvadingMode
		baseDirection = w.wrappedVector(nearestP, a.Pos)
	} else {
		a.mode = RoamingMode
	}

	finalDirection := baseDirection
	nearestPrey, _, foundPrey := w.nearestPreyEntityExcept(a.Pos, a)
	if foundPrey {
		// if a.mode == RoamingMode && rand.Float64() < clamp01(a.trait.roaming.val) {
		// 	peerHeading := nearestPrey.heading
		// 	if peerHeading.Length() == 0 {
		// 		peerHeading = randomUnitVector()
		// 	}
		// 	finalDirection = averageDirections(finalDirection, peerHeading)
		// }

		if a.mode == EvadingMode && rand.Float64() < clamp01(a.trait.selfishHerdChance.val) {
			towardPrey := w.wrappedVector(a.Pos, nearestPrey.Pos)
			finalDirection = averageDirections(finalDirection, towardPrey)
		}
	}

	if finalDirection.Length() == 0 {
		finalDirection = randomUnitVector()
	}
	a.heading = moveWithSpeed(finalDirection, 1)
	dpos := moveWithSpeed(finalDirection, preySpeed)
	w.move(&a.Pos, a.Char, dpos)
}

func (a *Prey) snapshot() CreatureSnapshot {
	return CreatureSnapshot{Pos: a.Pos, Char: a.Char, Kind: AnimalKind, Mode: a.mode, BaseColor: a.baseColor}
}

func (a *Prey) alive() bool {
	return !a.dead
}

func moveWithSpeed(direction u.Vector, speed float64) u.Vector {
	length := direction.Length()
	if length == 0 {
		return u.Vector{}
	}

	return u.Vector{X: direction.X / length * speed, Y: direction.Y / length * speed}
}

func randomUnitVector() u.Vector {
	angle := rand.Float64() * 2 * math.Pi
	return u.Vector{X: math.Cos(angle), Y: math.Sin(angle)}
}

func jitterDirection(direction u.Vector, maxDegrees float64) u.Vector {
	if direction.Length() == 0 {
		return randomUnitVector()
	}

	maxRadians := maxDegrees * math.Pi / 180.0
	delta := (rand.Float64()*2 - 1) * maxRadians
	angle := math.Atan2(direction.Y, direction.X) + delta
	return u.Vector{X: math.Cos(angle), Y: math.Sin(angle)}
}

func averageDirections(a, b u.Vector) u.Vector {
	return a.Add(b)
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
