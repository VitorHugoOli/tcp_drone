package genetic

import (
	"math"
	"math/rand"
	Model "tcp_drone/model"
)

type CrossOverFunc func(solution1, solution2 Model.Solution, params ...interface{}) []Model.Solution

// orderCrossover performs order crossover on two solutions
func orderCrossover(solution1, solution2 Model.Solution, params ...interface{}) Model.Solution {
	// init new solution
	newSolution := Model.Solution{
		Route: make([]int, len(solution1.Route)),
	}

	// init random indexes
	randomIndex1 := rand.Intn(len(solution1.Route))
	randomIndex2 := rand.Intn(len(solution1.Route))

	// init start and end indexes
	startIndex := int(math.Min(float64(randomIndex1), float64(randomIndex2)))
	endIndex := int(math.Max(float64(randomIndex1), float64(randomIndex2)))

	// init used
	used := make([]bool, len(solution1.Route))

	// init new solution route
	for i := startIndex; i <= endIndex; i++ {
		newSolution.Route[i] = solution1.Route[i]
		used[solution1.Route[i]] = true
	}

	// init new solution route
	for i := 0; i < len(solution1.Route); i++ {
		if i >= startIndex && i <= endIndex {
			continue
		}

		for j := 0; j < len(solution1.Route); j++ {
			if !used[solution2.Route[j]] {
				newSolution.Route[i] = solution2.Route[j]
				used[solution2.Route[j]] = true
				break
			}
		}
	}

	// return new solution
	return newSolution
}
