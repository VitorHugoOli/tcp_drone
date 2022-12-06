package genetic

import (
	"math/rand"
	Model "tcp_drone/model"
)

type SelectionFunc func(population []Model.Solution, params ...interface{}) Model.Solution

// tournamentSelection selects one solution from population using tournament method and returns its index
func tournamentSelection(population []Model.Solution, params ...interface{}) Model.Solution {
	// get tournament size
	var tournamentSize int
	if len(params) > 1 {
		tournamentSize = params[1].(int)
	} else {
		tournamentSize = 3
	}

	// get random solutions
	var solutions []Model.Solution
	for i := 0; i < tournamentSize; i++ {
		var s Model.Solution
		if rand.Intn(100) < 60 {
			s = population[rand.Intn(len(population)-len(population)/2)]
		} else {
			s = population[rand.Intn(len(population))]
		}
		s.Fitness()
		solutions = append(solutions, s)
	}

	// get best solution
	var bestSolution Model.Solution
	var bestFitness float64
	for i, solution := range solutions {
		if i == 0 || solution.RouteTime < bestFitness {
			bestSolution = solution
			bestFitness = solution.RouteTime
		}
	}

	return bestSolution
}

//// rouletteSelection selects two solutions from population using roulette method and returns their indexes
//func rouletteSelection(population []Model.Solution, params ...interface{}) []int {
//	// init fitness
//	fitness := params[0].([]float64)
//
//	// init selected
//	selected := make([]int, 0)
//
//	// select two solutions
//	for i := 0; i < 2; i++ {
//		// init random number
//		randomNumber := rand.Float64()
//
//		// init sum
//		sum := 0.0
//
//		// init index
//		index := 0
//
//		// find solution
//		for sum < randomNumber && index < len(population) {
//			sum += fitness[index]
//			index++
//		}
//
//		// add solution index to selected
//		selected = append(selected, index-1)
//	}
//
//	// return selected
//	return selected
//}
//
//// rankSelection selects two solutions from population using rank method and returns their indexes
//func rankSelection(population []Model.Solution, params ...interface{}) []int {
//	// init selected
//	selected := make([]int, 0)
//
//	// select two solutions
//	for i := 0; i < 2; i++ {
//		// init random number
//		randomNumber := rand.Float64()
//
//		// init sum
//		sum := 0.0
//
//		// init index
//		index := 0
//
//		// find solution
//		for sum < randomNumber && index < len(population) {
//			sum += 1 / float64(index+1)
//			index++
//		}
//
//		// add solution index to selected
//		selected = append(selected, index-1)
//	}
//
//	// return selected
//	return selected
//}
//
//// stochasticUniversalSelection selects two solutions from population using stochastic universal method and returns their indexes
//func stochasticUniversalSelection(population []Model.Solution, params ...interface{}) []int {
//	// init fitness
//	fitness := params[0].([]float64)
//
//	// init selected
//	selected := make([]int, 0)
//
//	// init sum
//	sum := 0.0
//
//	// calculate sum
//	for _, value := range fitness {
//		sum += value
//	}
//
//	// init distance
//	distance := sum / 2
//
//	// init start
//	start := rand.Float64() * distance
//
//	// select two solutions
//	for i := 0; i < 2; i++ {
//		// init sum
//		sum := 0.0
//
//		// init index
//		index := 0
//
//		// find solution
//		for sum < start && index < len(population) {
//			sum += fitness[index]
//			index++
//		}
//
//		// add solution index to selected
//		selected = append(selected, index-1)
//
//		// update start
//		start += distance
//	}
//
//	// return selected
//	return selected
//}
