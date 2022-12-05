package local_search

import (
	"tcp_drone/heuristics/builders"
	Model "tcp_drone/model"
)

func insertion(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) {
	builders.CreatingInitSolution(city, solution, builderHeuristic)
	tempSolution := Model.Solution{}
	holdBestTempSolution := Model.Solution{}

	tempSolution.City = solution.City
	holdBestTempSolution.City = solution.City

	tempSolution.Route = make([]int, solution.City.NodesLen+1)
	holdBestTempSolution.Route = make([]int, len(solution.Route))

	copy(holdBestTempSolution.Route, solution.Route)
	holdBestTempSolution.RouteTime = solution.RouteTime

	staticRoute := make([]int, len(solution.Route))
	copy(staticRoute, solution.Route)

	for i := 1; i < city.NodesLen; i++ {
		node := staticRoute[i]
		index, _ := solution.FindCustomerIndex(node)

		copy(tempSolution.Route, solution.Route)
		tempSolution.Route = append(tempSolution.Route[:index], tempSolution.Route[index+1:]...)

		for j := 1; j < city.NodesLen; j++ {

			//flush temp solution
			tempSolution.Route = append(tempSolution.Route[:j], append([]int{node}, tempSolution.Route[j:]...)...)
			tempSolution.Calculate()

			if tempSolution.RouteTime < holdBestTempSolution.RouteTime {
				copy(holdBestTempSolution.Route, tempSolution.Route)
				holdBestTempSolution.RouteTime = tempSolution.RouteTime
			}
			//remove node from temp solution
			tempSolution.Route = append(tempSolution.Route[:j], tempSolution.Route[j+1:]...)
			solution.Iterations++
		}

		tempSolution.Route = append(tempSolution.Route[:index], append([]int{node}, tempSolution.Route[index:]...)...)
		if holdBestTempSolution.RouteTime < solution.RouteTime {
			copy(solution.Route, holdBestTempSolution.Route)
			solution.RouteTime = holdBestTempSolution.RouteTime
			i = 0
		}

	}

}

func Insertion(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{LocalSearchAlgorithm: "Insertion"}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	defer solution.Timer()()
	insertion(city, solution, builderHeuristic)
	return solution, nil
}
