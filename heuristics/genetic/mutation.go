package genetic

import (
	"math/rand"
	"tcp_drone/heuristics/global_search"
	Model "tcp_drone/model"
)

type MutationFunc func(solution *Model.Solution)

// swapMutation performs swap mutation on a solution
func swapMutation(solution *Model.Solution) {

	global_search.PerturbationNodes(solution)

	// calculate new route time
	solution.Fitness()
}

// insertMutation performs insert mutation on a solution
func insertMutation(solution *Model.Solution) {
	// get random positions
	var (
		pos1 = rand.Intn(len(solution.Route))
		pos2 = rand.Intn(len(solution.Route))
	)

	// insert value
	solution.Route = append(solution.Route[:pos1], append([]int{solution.Route[pos2]}, solution.Route[pos1:pos2]...)...)

	// calculate new route time
	solution.Calculate()
}

// scrambleMutation performs scramble mutation on a solution
func scrambleMutation(solution *Model.Solution) {
	// get random positions
	var (
		pos1 = rand.Intn(len(solution.Route))
		pos2 = rand.Intn(len(solution.Route))
	)

	// scramble values
	solution.Route = append(solution.Route[:pos1], append(solution.Route[pos2:], solution.Route[pos1:pos2]...)...)

	// calculate new route time
	solution.Calculate()
}
