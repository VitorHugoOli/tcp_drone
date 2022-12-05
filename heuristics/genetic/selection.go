package genetic

import (
	"math/rand"
	Model "tcp_drone/model"
)

type SelectionFunc func(population []Model.Solution, params ...interface{}) []int

// tournamentSelection selects two random solutions from population and returns the best one index
func tournamentSelection(population []Model.Solution, params ...interface{}) []int {
	// init selected
	selected := make([]int, 0)

	// select two random solutions
	for i := 0; i < 2; i++ {
		// init random index
		randomIndex := rand.Intn(len(population))

		// init best solution
		bestSolution := population[randomIndex]

		// init best solution index
		bestSolutionIndex := randomIndex

		// find best solution
		for j := 0; j < 5; j++ {
			// init random index
			randomIndex := rand.Intn(len(population))

			// compare solutions
			if population[randomIndex].RouteTime < bestSolution.RouteTime {
				bestSolution = population[randomIndex]
				bestSolutionIndex = randomIndex
			}
		}

		// add best solution index to selected
		selected = append(selected, bestSolutionIndex)
	}

	// return selected
	return selected
}

// rouletteSelection selects two solutions from population using roulette method and returns their indexes
func rouletteSelection(population []Model.Solution, params ...interface{}) []int {
	// init fitness
	fitness := params[0].([]float64)

	// init selected
	selected := make([]int, 0)

	// select two solutions
	for i := 0; i < 2; i++ {
		// init random number
		randomNumber := rand.Float64()

		// init sum
		sum := 0.0

		// init index
		index := 0

		// find solution
		for sum < randomNumber && index < len(population) {
			sum += fitness[index]
			index++
		}

		// add solution index to selected
		selected = append(selected, index-1)
	}

	// return selected
	return selected
}

// rankSelection selects two solutions from population using rank method and returns their indexes
func rankSelection(population []Model.Solution, params ...interface{}) []int {
	// init selected
	selected := make([]int, 0)

	// select two solutions
	for i := 0; i < 2; i++ {
		// init random number
		randomNumber := rand.Float64()

		// init sum
		sum := 0.0

		// init index
		index := 0

		// find solution
		for sum < randomNumber && index < len(population) {
			sum += 1 / float64(index+1)
			index++
		}

		// add solution index to selected
		selected = append(selected, index-1)
	}

	// return selected
	return selected
}

// stochasticUniversalSelection selects two solutions from population using stochastic universal method and returns their indexes
func stochasticUniversalSelection(population []Model.Solution, params ...interface{}) []int {
	// init fitness
	fitness := params[0].([]float64)

	// init selected
	selected := make([]int, 0)

	// init sum
	sum := 0.0

	// calculate sum
	for _, value := range fitness {
		sum += value
	}

	// init distance
	distance := sum / 2

	// init start
	start := rand.Float64() * distance

	// select two solutions
	for i := 0; i < 2; i++ {
		// init sum
		sum := 0.0

		// init index
		index := 0

		// find solution
		for sum < start && index < len(population) {
			sum += fitness[index]
			index++
		}

		// add solution index to selected
		selected = append(selected, index-1)

		// update start
		start += distance
	}

	// return selected
	return selected
}
