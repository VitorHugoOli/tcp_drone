package local_search

import (
	"tcp_drone/heuristics/builders"
	Model "tcp_drone/model"
)

func threeOpt(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) {
	builders.CreatingInitSolution(city, solution, builderHeuristic)

	bestTime := solution.RouteTime

	for i := 1; i < city.NodesLen; i++ {
		for j := i + 1; j < city.NodesLen; j++ {
			for k := j + 1; k < city.NodesLen; k++ {
				if solution.Route[i] != solution.Route[j] && solution.Route[i] != solution.Route[k] && solution.Route[j] != solution.Route[k] {
					solution.Route[i], solution.Route[j], solution.Route[k] = solution.Route[k], solution.Route[i], solution.Route[j]
					solution.Calculate()
					if solution.RouteTime < bestTime {
						bestTime = solution.RouteTime
					} else {
						solution.Route[i], solution.Route[j], solution.Route[k] = solution.Route[j], solution.Route[k], solution.Route[i]
					}
				}
				solution.Iterations++
			}
		}
	}
}

func ThreeOpt(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{LocalSearchAlgorithm: "3-opt"}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	defer solution.Timer()()
	threeOpt(city, solution, builderHeuristic)
	return solution, nil
}
