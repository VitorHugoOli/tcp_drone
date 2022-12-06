package genetic

import (
	"math/rand"
	"tcp_drone/heuristics/global_search"
	Model "tcp_drone/model"
)

type MutationFunc func(solution *Model.Solution)

// swapMutation performs swap mutation on a solution
func swapMutation(solution *Model.Solution) {
	// get random positions
	//var (
	//	pos1 = rand.Intn(len(solution.Route)-2) + 1
	//	pos2 = rand.Intn(len(solution.Route)-2) + 1
	//)
	//
	//for pos1 == 0 || pos2 == 0 || pos1 == len(solution.Route)-1 || pos2 == len(solution.Route)-1 || solution.Route[pos1] == solution.Route[pos2] {
	//	pos1 = rand.Intn(len(solution.Route)-2) + 1
	//	pos2 = rand.Intn(len(solution.Route)-2) + 1
	//}
	//
	//// swap values
	//solution.Route[pos1], solution.Route[pos2] = solution.Route[pos2], solution.Route[pos1]
	//
	//for pos1 == pos2 || pos1 == 0 || pos2 == 0 || pos1 == len(solution.Route)-1 || pos2 == len(solution.Route)-1 {
	//	pos1 = rand.Intn(len(solution.Route)-2) + 1
	//	pos2 = rand.Intn(len(solution.Route)-2) + 1
	//}
	//
	//// change sign of values
	//solution.Route[pos1] = -solution.Route[pos1]
	//solution.Route[pos2] = -solution.Route[pos2]

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
