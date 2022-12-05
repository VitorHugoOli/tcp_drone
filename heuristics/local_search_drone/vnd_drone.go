package local_search

import (
	"tcp_drone/heuristics/builders"
	"tcp_drone/heuristics/builders_drone"
	"tcp_drone/heuristics/local_search"
	Model "tcp_drone/model"
)

func vndDrone(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) {
	if solution.RouteTime == 0 {
		if builderHeuristic == nil {
			_, _ = builders_drone.DroneInitialSolution(city, solution, nil)
		} else {
			_, _ = builderHeuristic(city, solution)
		}
	}

	localSearchers := []func(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error){
		local_search.TwoOpt,
		local_search.ThreeOpt,
		local_search.Insertion, // greedy insertion d=1
		Inversion,
	}

	local_search.VndLogic(city, solution, localSearchers)

}

func VndDrone(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	solution.LocalSearchAlgorithm = "VND Drone"
	defer solution.Timer()()
	vndDrone(city, solution, builderHeuristic)
	return solution, nil
}
