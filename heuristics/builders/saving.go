package builders

import (
	Model "tcp_drone/model"
)

func savings(city *Model.City, solution *Model.Solution) {
	initSolution := make([]int, (city.NodesLen*2)-1)
	count := 1
	for i := 1; i < len(initSolution); i = i + 2 {
		initSolution[i] = count
		count++
	}
}

func Savings(city *Model.City) (*Model.Solution, error) {
	solution := &Model.Solution{BuilderAlgorithm: "Savings"}
	defer solution.Timer()()
	e := solution.Init(*city)
	if e != nil {
		return nil, e
	}
	savings(city, solution)
	return solution, nil
}
