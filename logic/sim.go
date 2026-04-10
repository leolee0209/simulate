package logic

var defaultWorld = NewWorld()

func InitSimulation() {
	defaultWorld = NewWorld()
}

func Step() {
	defaultWorld.Step()
}

func Snapshot() []CreatureSnapshot {
	return defaultWorld.Snapshot()
}

func AddAnimals(num int) {
	defaultWorld.AddAnimals(num)
}

func AddAnimal(a Prey) {
	defaultWorld.AddAnimal(a)
}

func AddPredator(p Predator) {
	defaultWorld.AddPredator(p)
}

func AddPredators(num int) {
	defaultWorld.AddPredators(num)
}

func Dimensions() (float64, float64) {
	return defaultWorld.Size()
}

func Tick() uint64 {
	return defaultWorld.Tick()
}

func Generation() int {
	return defaultWorld.Generation()
}

func LastCompletedGeneration() int {
	return defaultWorld.LastCompletedGeneration()
}

func LastCompletedGenerationStats() (float64, int) {
	return defaultWorld.LastCompletedGenerationStats()
}

func CurrentGenerationStats() (float64, int) {
	return defaultWorld.CurrentGenerationStats()
}

func GenerationHistory() []GenerationStatsPoint {
	return defaultWorld.GenerationHistory()
}
