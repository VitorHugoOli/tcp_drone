package builders_drone

import (
	"tcp_drone/heuristics/builders"
	"tcp_drone/heuristics/local_search"
	Model "tcp_drone/model"
)

func DroneInitialSolution(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{LocalSearchAlgorithm: "Drone"}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	defer solution.Timer()()
	err := solution.Init(*city)
	if err != nil {
		return nil, err
	}
	droneInitialSolution(city, solution)
	return solution, nil
}

func droneInitialSolution(city *Model.City, solution *Model.Solution) {
	if solution.RouteTime == 0 {
		_, _ = local_search.Vnd(city, solution, nil)
	}

	subSolution := Model.SubSolution{0, false, 0, 0, 0, 0}

	//dronesAvv := make([][]int, 0)
	//for i, _ := range solution.City.DroneTimeDistance {
	//	for j := i + 1; j < len(solution.City.DroneTimeDistance[i]); j++ {
	//		if solution.City.DroneTimeDistance[i][j] <= solution.City.DroneSettings.BatteryTime {
	//			print("Drone ", i, " can deliver package ", j, " in ", solution.City.DroneTimeDistance[i][j], " minutes\n")
	//			dronesAvv = append(dronesAvv, []int{i, j})
	//		}
	//	}
	//}

	var excludeNodes []int
	for {
		for _, customer := range city.GetEligibleDroneCustomers(&excludeNodes) {
			index, _ := solution.FindCustomerIndex(customer.Node)
			subRoutes := solution.GetSubRoutes()
			subSolution.Savings = solution.CalcSavings(index, subRoutes)
			for _, subRoute := range subRoutes {
				if !subRoute.OnlyDriver {
					solution.CalcCostDriver(subRoute, customer.Node, &subSolution)
				} else {
					solution.CalcCostDrone(subRoute, customer.Node, &subSolution)
				}
			}
		}
		if subSolution.MaxSaving > 0 {
			solution.SwapCustomer(subSolution.Customer, subSolution.Departure)
			if subSolution.IsDrone {
				excludeNodes = append(excludeNodes, subSolution.Customer*-1)
				excludeNodes = append(excludeNodes, subSolution.Departure)
				excludeNodes = append(excludeNodes, subSolution.Arrival)
			} else {
				excludeNodes = append(excludeNodes, subSolution.Customer)
			}
			subSolution = Model.SubSolution{0, false, 0, 0, 0, 0}
		} else {
			break
		}
	}
}
