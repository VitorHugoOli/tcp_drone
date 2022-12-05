package builders

import (
	Model "tcp_drone/model"
)

type BuilderHeuristic func(city *Model.City, solution *Model.Solution) (*Model.Solution, error)

func CreatingInitSolution(city *Model.City, solution *Model.Solution, builderHeuristic BuilderHeuristic) {
	if solution.RouteTime == 0 {
		if builderHeuristic == nil {
			_, _ = NearestNeighbor(city, solution)
		} else {
			_, _ = builderHeuristic(city, solution)
		}
	}
}
