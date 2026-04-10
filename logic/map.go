package logic

import (
	"math"
	"math/rand"
	"simulate/util"
)

const Y_SIZE = 100.0
const X_SIZE = 100.0

var cycleRatio float64 = 0.5

type Creature interface {
	act(w *World)
	snapshot() CreatureSnapshot
	alive() bool
}

type CreatureKind int

const (
	AnimalKind CreatureKind = iota
	PredatorKind
)

type BehaviorMode int

const (
	IdleMode BehaviorMode = iota
	RoamingMode
	EvadingMode
	ChasingMode
	RestingMode
)

type ColorRGB struct {
	R uint8
	G uint8
	B uint8
}

type CreatureSnapshot struct {
	Pos       util.Position
	Char      byte
	Kind      CreatureKind
	Mode      BehaviorMode
	BaseColor ColorRGB
}

type GenerationStatsPoint struct {
	Generation           int
	AvgSelfishHerdChance float64
	SurvivorCount        int
}

type World struct {
	width                             float64
	height                            float64
	creatures                         []Creature
	tick                              uint64
	initialPreyCount                  int
	initialPredatorCount              int
	generation                        int
	lastCompletedGen                  int
	lastCompletedAvgSelfishHerdChance float64
	lastCompletedSurvivorCount        int
	initialBaselineSet                bool
	initialBaseline                   GenerationStatsPoint
	completedHistory                  []GenerationStatsPoint
}

func NewWorld() *World {
	w := &World{width: X_SIZE, height: Y_SIZE}
	w.InitGrid()
	return w
}

func (w *World) AddAnimals(num int) {
	for range num {
		//get random ascii char, reserve p for predator
		var char byte
		for {
			char = byte(rand.Intn(94) + 33)
			if char != '.' && char != 'P' && char != 'p' {
				break
			}
		}

		x := rand.Float64() * w.width
		y := rand.Float64() * w.height

		//Add
		w.AddAnimal(Prey{Pos: util.Position{X: x, Y: y}, Char: char})
	}
}

func (w *World) AddAnimal(a Prey) {
	a.trait.vision.init()
	a.trait.roaming.init()
	a.trait.selfishHerdChance.init()
	if w.generation == 0 {
		w.initialPreyCount++
	}
	w.placePrey(a)
}

func (w *World) AddPredator(p Predator) {
	p.vision.init()
	if w.generation == 0 {
		w.initialPredatorCount++
	}
	w.placePredator(p)
}

func (w *World) AddPredators(num int) {
	for range num {
		x := rand.Float64() * w.width
		y := rand.Float64() * w.height

		w.AddPredator(Predator{Pos: util.Position{X: x, Y: y}})
	}
}

func (w *World) Step() {
	if w.generation == 0 && !w.initialBaselineSet {
		avg, survivors := w.CurrentGenerationStats()
		w.initialBaseline = GenerationStatsPoint{
			Generation:           0,
			AvgSelfishHerdChance: avg,
			SurvivorCount:        survivors,
		}
		w.initialBaselineSet = true
	}

	for i := range w.creatures {
		if w.creatures[i].alive() {
			w.creatures[i].act(w)
		}
	}
	w.compactCreatures()
	if w.shouldRestartGeneration() {
		w.restartGeneration()
		return
	}
	w.tick++
}

func (w *World) Snapshot() []CreatureSnapshot {
	res := make([]CreatureSnapshot, 0, len(w.creatures))
	for i := range w.creatures {
		res = append(res, w.creatures[i].snapshot())
	}
	return res
}

func (w *World) Tick() uint64 {
	return w.tick
}

func (w *World) Generation() int {
	return w.generation
}

func (w *World) LastCompletedGeneration() int {
	return w.lastCompletedGen
}

func (w *World) LastCompletedGenerationStats() (float64, int) {
	return w.lastCompletedAvgSelfishHerdChance, w.lastCompletedSurvivorCount
}

func (w *World) CurrentGenerationStats() (float64, int) {
	current := w.survivingPrey()
	return averageSurvivorSelfishHerdChance(current), len(current)
}

func (w *World) InitGrid() {
	w.creatures = nil
	w.tick = 0
	w.initialPreyCount = 0
	w.initialPredatorCount = 0
	w.generation = 0
	w.lastCompletedGen = 0
	w.lastCompletedAvgSelfishHerdChance = 0
	w.lastCompletedSurvivorCount = 0
	w.initialBaselineSet = false
	w.initialBaseline = GenerationStatsPoint{}
	w.completedHistory = nil
}

func (w *World) GenerationHistory() []GenerationStatsPoint {
	result := make([]GenerationStatsPoint, 0, len(w.completedHistory)+1)
	if w.initialBaselineSet {
		result = append(result, w.initialBaseline)
	} else {
		avg, survivors := w.CurrentGenerationStats()
		result = append(result, GenerationStatsPoint{
			Generation:           0,
			AvgSelfishHerdChance: avg,
			SurvivorCount:        survivors,
		})
	}

	result = append(result, w.completedHistory...)
	return result
}

func (w *World) move(pos *util.Position, char byte, dpos util.Vector) bool {
	_ = char
	newPos := *pos
	newPos = newPos.Add(dpos)

	newPos.X = math.Mod(newPos.X, w.width)
	newPos.Y = math.Mod(newPos.Y, w.height)
	if newPos.X < 0 {
		newPos.X += w.width
	}
	if newPos.Y < 0 {
		newPos.Y += w.height
	}

	pos.X = newPos.X
	pos.Y = newPos.Y
	return true
}

func wrappedDelta(from, to, extent float64) float64 {
	delta := to - from
	if delta > extent/2 {
		delta -= extent
	} else if delta < -extent/2 {
		delta += extent
	}

	return delta
}

func (w *World) wrappedVector(from, to util.Position) util.Vector {
	return util.Vector{
		X: wrappedDelta(from.X, to.X, w.width),
		Y: wrappedDelta(from.Y, to.Y, w.height),
	}
}

func (w *World) wrappedDistance(from, to util.Position) float64 {
	return w.wrappedVector(from, to).Length()
}

func (w *World) nearestPredator(from util.Position) (util.Position, bool, float64) {
	bestDistance := math.MaxFloat64
	nearest := util.Position{X: -1, Y: -1}

	for i := range w.creatures {
		if !w.creatures[i].alive() {
			continue
		}
		s := w.creatures[i].snapshot()
		if s.Kind != PredatorKind {
			continue
		}

		distance := w.wrappedDistance(from, s.Pos)
		if distance < bestDistance {
			bestDistance = distance
			nearest = s.Pos
		}
	}

	return nearest, !nearest.Equal(util.Position{X: -1, Y: -1}), bestDistance
}

func (w *World) nearestPrey(from util.Position) (util.Position, bool, float64) {
	return w.nearestPreyExcept(from, util.Position{X: math.MaxFloat64, Y: math.MaxFloat64})
}

func (w *World) nearestPreyExcept(from util.Position, exclude util.Position) (util.Position, bool, float64) {
	prey, bestDistance, found := w.nearestPreyEntityExcept(from, nil)
	if !found {
		return util.Position{X: -1, Y: -1}, false, math.MaxFloat64
	}

	if prey.Pos.Equal(exclude) {
		prey, bestDistance, found = w.nearestPreyEntityExcept(from, prey)
		if !found {
			return util.Position{X: -1, Y: -1}, false, math.MaxFloat64
		}
	}

	return prey.Pos, true, bestDistance
}

func (w *World) nearestPreyEntityExcept(from util.Position, exclude *Prey) (*Prey, float64, bool) {
	bestDistance := math.MaxFloat64
	var nearest *Prey

	for i := range w.creatures {
		if !w.creatures[i].alive() {
			continue
		}
		prey, ok := w.creatures[i].(*Prey)
		if !ok {
			continue
		}
		if exclude != nil && prey == exclude {
			continue
		}

		distance := w.wrappedDistance(from, prey.Pos)
		if distance < bestDistance {
			bestDistance = distance
			nearest = prey
		}
	}

	if nearest == nil {
		return nil, math.MaxFloat64, false
	}

	return nearest, bestDistance, true
}

func (w *World) Size() (float64, float64) {
	return w.width, w.height
}

func (w *World) preyCount() int {
	count := 0
	for i := range w.creatures {
		if !w.creatures[i].alive() {
			continue
		}
		if w.creatures[i].snapshot().Kind == AnimalKind {
			count++
		}
	}
	return count
}

func (w *World) shouldRestartGeneration() bool {
	if w.initialPreyCount == 0 {
		return false
	}

	return w.preyCount() <= int(float64(w.initialPreyCount)*cycleRatio)
}

func (w *World) eatPreyAt(pos util.Position, captureRadius float64) bool {
	for i := range w.creatures {
		if !w.creatures[i].alive() {
			continue
		}
		prey, ok := w.creatures[i].(*Prey)
		if !ok {
			continue
		}

		distance := w.wrappedDistance(pos, prey.Pos)
		if distance > captureRadius {
			continue
		}

		prey.dead = true
		return true
	}

	return false
}

func (w *World) compactCreatures() {
	kept := w.creatures[:0]
	for i := range w.creatures {
		if w.creatures[i].alive() {
			kept = append(kept, w.creatures[i])
		}
	}
	w.creatures = kept
}

func (w *World) restartGeneration() {
	survivors := w.survivingPrey()
	w.lastCompletedGen = w.generation + 1
	w.lastCompletedAvgSelfishHerdChance = averageSurvivorSelfishHerdChance(survivors)
	w.lastCompletedSurvivorCount = len(survivors)
	w.completedHistory = append(w.completedHistory, GenerationStatsPoint{
		Generation:           w.lastCompletedGen,
		AvgSelfishHerdChance: w.lastCompletedAvgSelfishHerdChance,
		SurvivorCount:        w.lastCompletedSurvivorCount,
	})
	w.creatures = nil
	w.tick = 0
	w.generation++

	if len(survivors) == 0 {
		w.spawnRandomGeneration()
		return
	}

	newPrey := w.nextGenerationPrey(survivors)
	for i := range newPrey {
		w.placePrey(newPrey[i])
	}

	for i := 0; i < w.initialPredatorCount; i++ {
		predator := Predator{Pos: util.Position{X: rand.Float64() * w.width, Y: rand.Float64() * w.height}}
		predator.vision.init()
		w.placePredator(predator)
	}
}

func (w *World) survivingPrey() []Prey {
	res := make([]Prey, 0)
	for i := range w.creatures {
		if !w.creatures[i].alive() {
			continue
		}
		prey, ok := w.creatures[i].(*Prey)
		if !ok {
			continue
		}
		res = append(res, *prey)
	}
	return res
}

func (w *World) nextGenerationPrey(survivors []Prey) []Prey {
	targetCount := w.initialPreyCount
	if targetCount <= 0 {
		targetCount = len(survivors) * 2
	}

	children := make([]Prey, 0, targetCount)
	for i := range survivors {
		children = append(children, w.makeChildPrey(survivors[i], -1))
		children = append(children, w.makeChildPrey(survivors[i], 1))
	}

	for len(children) < targetCount {
		parent := survivors[rand.Intn(len(survivors))]
		children = append(children, w.makeChildPrey(parent, 1))
	}

	if len(children) > targetCount {
		children = children[:targetCount]
	}

	return children
}

func (w *World) makeChildPrey(parent Prey, direction float64) Prey {
	child := parent
	child.dead = false
	child.mode = RoamingMode
	child.Pos.X = math.Mod(parent.Pos.X+direction*0.7+rand.Float64()*0.4-0.2, w.width)
	child.Pos.Y = math.Mod(parent.Pos.Y+rand.Float64()*0.4-0.2, w.height)
	if child.Pos.X < 0 {
		child.Pos.X += w.width
	}
	if child.Pos.Y < 0 {
		child.Pos.Y += w.height
	}
	child.trait.roaming.val = clampFloat(parent.trait.roaming.val+rand.Float64()*0.2-0.1, child.trait.roaming.min, child.trait.roaming.max)
	child.trait.selfishHerdChance.val = clampFloat(parent.trait.selfishHerdChance.val+rand.Float64()*0.1-0.05, child.trait.selfishHerdChance.min, child.trait.selfishHerdChance.max)
	return child
}

func (w *World) spawnRandomGeneration() {
	for i := 0; i < w.initialPreyCount; i++ {
		x := rand.Float64() * w.width
		y := rand.Float64() * w.height
		w.AddAnimal(Prey{Pos: util.Position{X: x, Y: y}, Char: byte(rand.Intn(94) + 33)})
	}

	for i := 0; i < w.initialPredatorCount; i++ {
		x := rand.Float64() * w.width
		y := rand.Float64() * w.height
		w.AddPredator(Predator{Pos: util.Position{X: x, Y: y}})
	}
}

func (w *World) placePrey(a Prey) {
	a.mode = RoamingMode
	a.baseColor = ColorRGB{R: 70, G: 120, B: 255}
	a.dead = false
	if a.heading.Length() == 0 {
		a.heading = randomUnitVector()
	}
	w.creatures = append(w.creatures, &a)
	w.move(&a.Pos, a.Char, util.Vector{})
}

func (w *World) placePredator(p Predator) {
	p.mode = RestingMode
	p.baseColor = ColorRGB{R: 220, G: 50, B: 50}
	p.restTicksLeft = 0
	p.dead = false
	w.creatures = append(w.creatures, &p)
	w.move(&p.Pos, 'P', util.Vector{})
}

func clampFloat(value, minValue, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func averageSurvivorSelfishHerdChance(survivors []Prey) float64 {
	if len(survivors) == 0 {
		return 0
	}

	total := 0.0
	for i := range survivors {
		total += survivors[i].trait.selfishHerdChance.val
	}

	return total / float64(len(survivors))
}

func SetCycle(ratio float64) {
	cycleRatio = ratio
}
