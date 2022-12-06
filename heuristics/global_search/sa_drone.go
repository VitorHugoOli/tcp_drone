package global_search

import (
	"math"
	"math/rand"
	"tcp_drone/heuristics/builders"
	local_search_drone "tcp_drone/heuristics/local_search_drone"
	Model "tcp_drone/model"
	"time"
)

type SimulatedSettings struct {
	MaxTemp   float64
	MinTemp   float64
	maxIter   int
	alpha     float64
	successes int
}

func defaultSimulatedSettings() SimulatedSettings {
	return SimulatedSettings{
		MaxTemp:   20000,
		MinTemp:   0.1,
		maxIter:   200,
		alpha:     0.9,
		successes: 10,
	}
}

// simulated Annealing
func SimulatedAnnealing(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	solution.GlobalSearchAlgorithm = "Simulated Annealing"
	defer solution.Timer()()
	simulatedAnnealing(city, solution, builderHeuristic)
	return solution, nil
}

// generate a reandom solution given a current solution
func randomSolution(city *Model.City, solution *Model.Solution) *Model.Solution {
	newSolution := &Model.Solution{}
	_ = newSolution.Init(*city)
	newSolution.Route = make([]int, len(solution.Route))
	copy(newSolution.Route, solution.Route)
	newSolution.RouteTime = solution.RouteTime

	// calculate new route time
	PerturbationNodes(newSolution)
	//PermutationSignal(newSolution)
	//perturbationSubGroups(newSolution)

	newSolution.Calculate()

	return newSolution
}

func PerturbationNodes(newSolution *Model.Solution) {
	rand.Seed(time.Now().UnixNano())
	// sort a number of nodes to be swapped
	var dice int
	for {
		dice = rand.Intn(6)
		if dice >= 2 && dice%2 == 0 {
			break
		}
	}

	// get the nodes to be swapped
	arr := make([]int, dice)
	for i := 0; i < dice; i++ {
		arr[i] = rand.Intn(len(newSolution.Route)-2) + 1
	}

	// swap the nodes
	for i := 0; i < len(arr)/2; i++ {
		newSolution.Route[arr[i]], newSolution.Route[arr[len(arr)-1-i]] = newSolution.Route[arr[len(arr)-1-i]], newSolution.Route[arr[i]]
	}

	// 90% chance to swap the signal
	for i := 0; i < len(arr); i++ {
		if rand.Intn(10) < 5 {
			newSolution.Route[arr[i]] = -newSolution.Route[arr[i]]
		}
	}
}

func PermutationSignal(solution *Model.Solution) {
	var dice int
	// role the dice again
	for {
		dice = rand.Intn(len(solution.Route))
		if dice > 2 && dice%2 == 0 {
			break
		}
	}

	// get the nodes to be change signal
	arr := make([]int, dice)
	for i := 0; i < dice; i++ {
		arr[i] = rand.Intn(len(solution.Route)-2) + 1
	}

	// change signal
	for i := 0; i < len(arr); i++ {
		solution.Route[arr[i]] = -solution.Route[arr[i]]
	}
}

func perturbationSubGroups(solution *Model.Solution) {
	subGroups := solution.GetSubRoutes()

	//get two random subroutes with drone
	rand.Seed(time.Now().UnixNano())
	var i, j int

	for {
		i = rand.Intn(len(subGroups))
		j = rand.Intn(len(subGroups))

		if i != j && subGroups[i].OnlyDriver && subGroups[j].OnlyDriver && i != 0 && j != 0 && i != len(subGroups)-1 && j != len(subGroups)-1 {
			break
		}
	}

	subGroups[i].Route, subGroups[j].Route = subGroups[j].Route, subGroups[i].Route
	//print route
	solution.Route = Model.AppendSubRoute(subGroups)

}

func acceptanceProbability(currentRouteTime float64, newRouteTime float64, temp float64) bool {
	if newRouteTime < currentRouteTime {
		return true
	}
	if newRouteTime == math.MaxFloat64 {
		return false
	}
	rand.Seed(time.Now().UnixNano())
	return math.Exp((currentRouteTime-newRouteTime)/temp) > rand.Float64()
}

func simulatedAnnealing(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) {
	if solution.RouteTime == 0 {
		if builderHeuristic == nil {
			_, _ = local_search_drone.VndDrone(city, solution, nil)
		} else {
			_, _ = builderHeuristic(city, solution)
		}
	}

	settings := defaultSimulatedSettings()
	temp := settings.MaxTemp
	iter := 0

	tempSolution := &Model.Solution{}
	_ = tempSolution.Init(*city)
	tempSolution.Route = make([]int, len(solution.Route))
	copy(tempSolution.Route, solution.Route)
	tempSolution.RouteTime = solution.RouteTime

	for {
		for {
			newSolution := randomSolution(city, tempSolution)
			if acceptanceProbability(tempSolution.RouteTime, newSolution.RouteTime, temp) {
				tempSolution.Route = newSolution.Route
				tempSolution.RouteTime = newSolution.RouteTime
			}
			if iter >= settings.maxIter {
				break
			}
			iter++
			tempSolution.Iterations++
		}
		_, _ = local_search_drone.VndDrone(city, tempSolution, nil)
		temp *= settings.alpha

		if tempSolution.RouteTime < solution.RouteTime {
			solution.Route = make([]int, len(tempSolution.Route))
			copy(solution.Route, tempSolution.Route)
			solution.RouteTime = tempSolution.RouteTime
		}

		if temp < settings.MinTemp {
			break
		}
		iter = 0
	}
}
