package builders

import (
	Model "tcp_drone/model"
)

// NearestNeighbourHeuristic is a heuristic that finds the nearest neighbour
// of the current node and adds it to the route.

func nearestNeighbour(city *Model.City, solution *Model.Solution) {
	for i := 1; i < city.NodesLen; i++ {
		NearestNeighbourLogic(city, solution, i)
	}
}

func NearestNeighbourLogic(city *Model.City, solution *Model.Solution, i int) {
	currentRoute := solution.Route[:i-1]
	nearest := city.GetNearestNodeDriver(solution.Route[i-1], &currentRoute)
	solution.Route[i] = nearest.Node
	solution.Iterations++
}

func NearestNeighbor(city *Model.City, solution *Model.Solution) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	solution.BuilderAlgorithm = "Nearest Neighbour"
	defer solution.Timer()()

	nearestNeighbour(city, solution)
	return solution, nil
}
