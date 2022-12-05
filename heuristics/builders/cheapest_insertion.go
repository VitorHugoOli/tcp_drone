package builders

import (
	Model "tcp_drone/model"
	"tcp_drone/utils"
)

type CIModel struct {
	index int
	cost  float64
}

func cheapestInsertion(city *Model.City, solution *Model.Solution) {
	for i := 1; i <= 3; i++ {
		NearestNeighbourLogic(city, solution, i)
	}

	customers := make([]int, city.NodesLen)
	for i := range customers {
		customers[i] = i
	}

	remainingCustomers := utils.Difference(customers, solution.Route)

	for i := 4; i < city.NodesLen; i++ {
		bestCustomer := CIModel{cost: 1000000000000000000}
		for _, customer := range remainingCustomers {
			for j := 0; j < i; j++ {
				cost := city.GetDriverTimeDistance(solution.Route[j], customer) + city.GetDriverTimeDistance(customer, solution.Route[j+1]) - city.GetDriverTimeDistance(solution.Route[j], solution.Route[j+1])
				if cost < bestCustomer.cost {
					bestCustomer.cost = cost
					bestCustomer.index = customer
					solution.Iterations++
				}
			}
		}
		if len(remainingCustomers) >= 1 {
			remainingCustomers = utils.Remove(remainingCustomers, bestCustomer.index)
		}
		solution.Route[i] = bestCustomer.index
	}
}

func CheapestInsertion(city *Model.City, solution *Model.Solution) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{BuilderAlgorithm: "Nearest Neighbour"}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	defer solution.Timer()()
	e := solution.Init(*city)
	if e != nil {
		return nil, e
	}
	cheapestInsertion(city, solution)
	return solution, nil
}
