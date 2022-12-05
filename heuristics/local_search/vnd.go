package local_search

import (
	"tcp_drone/heuristics/builders"
	Model "tcp_drone/model"
)

func vnd(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) {
	builders.CreatingInitSolution(city, solution, builderHeuristic)
	solution.CheckDuplicateItemsRoute()
	solution.CheckAllItemsIndRoute()

	localSearchers := []func(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error){
		Insertion,
		TwoOpt,
		ThreeOpt,
	}

	VndLogic(city, solution, localSearchers)
}

func VndLogic(city *Model.City, solution *Model.Solution, localSearchers []func(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error)) {
	bestTime := solution.RouteTime
	k := 0
	improves := 0
	for {
		solution, _ = localSearchers[k](city, solution, nil)
		if solution.RouteTime < bestTime {
			bestTime = solution.RouteTime
			k = 0
			improves++
		} else {
			k++
		}
		if k == len(localSearchers) {
			k = 0
			if solution.RouteTime == bestTime && improves == 0 {
				break
			}
			improves = 0
		}
	}
}

func Vnd(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{LocalSearchAlgorithm: "VND"}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	defer solution.Timer()()
	vnd(city, solution, builderHeuristic)
	return solution, nil
}
