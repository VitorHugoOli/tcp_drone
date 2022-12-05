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
		MaxTemp:   5000,
		MinTemp:   0.1,
		maxIter:   500,
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

	//subroutes := newSolution.GetSubRoutes()

	//rand.Seed(time.Now().UnixNano())
	//i := rand.Intn(len(subroutes)-2) + 1
	//j := rand.Intn(len(subroutes)-2) + 1
	//subroutes[i], subroutes[j] = subroutes[j], subroutes[i]
	//
	//// join subroutesand remove duplicates
	//newSolution.Route = []int{}
	//for _, subroute := range subroutes {
	//	newSolution.Route = append(newSolution.Route, subroute.Route...)
	//}
	//
	//newSolution.Route = removeDuplicates(newSolution.Route)

	// shuffle the positives nodes at routes
	for {
		rand.Seed(time.Now().UnixNano())
		i := rand.Intn(len(newSolution.Route)-2) + 1
		j := rand.Intn(len(newSolution.Route)-2) + 1
		if newSolution.Route[i] > 0 && newSolution.Route[j] > 0 {
			newSolution.Route[i], newSolution.Route[j] = newSolution.Route[j], newSolution.Route[i]
			break
		}
	}

	// calculate new route time
	newSolution.Calculate()

	return newSolution
}

func pertubation(newSolution *Model.Solution) {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(newSolution.Route)-2) + 1
	j := rand.Intn(len(newSolution.Route)-2) + 1

	newSolution.Route[i], newSolution.Route[j] = newSolution.Route[j], newSolution.Route[i]
}

func acceptanceProbability(currentRouteTime float64, newRouteTime float64, temp float64) bool {
	if newRouteTime < currentRouteTime {
		return true
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
		//_, _ = local_search_drone.VndDrone(city, tempSolution, nil)
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
