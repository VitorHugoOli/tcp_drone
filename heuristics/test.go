package heuristics

import (
	"fmt"
	"tcp_drone/heuristics/builders"
	"tcp_drone/heuristics/builders_drone"
	"tcp_drone/heuristics/global_search"
	"tcp_drone/heuristics/local_search"
	localDrone "tcp_drone/heuristics/local_search_drone"
	Model "tcp_drone/model"
)

func Test() {
	//Testing nearest neighbour heuristic
	c := &Model.City{}
	//c.Init("tsplib_generated_instances/berlin52/", nil, nil)
	c.Init("tsplib_generated_instances/eil51/", nil, nil)
	//c.Init("tsplib_generated_instances/bier127/", nil, nil)

	//nearest(c)
	//nearestModify(c)
	//cheapest(c)
	//twoOpt(c)
	//insertion(c)
	//threeOpt(c)
	//vnd(c)
	//drone(c)
	//vndDrone(c)
	simulateAnnealing(c)
}

func simulateAnnealing(c *Model.City) {
	s, err := global_search.SimulatedAnnealing(c, nil, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return
	}
	s.DebugPrint()
}

func drone(c *Model.City) {
	s, _ := builders_drone.DroneInitialSolution(c, nil, nil)
	s.DebugPrint()
}

func vndDrone(c *Model.City) {
	s, err := localDrone.VndDrone(c, nil, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return
	}
	s.DebugPrint()
}

func vnd(c *Model.City) {
	s, err := local_search.Vnd(c, nil, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return
	}
	s.DebugPrint()
}

func twoOpt(c *Model.City) {
	s, err := local_search.TwoOpt(c, nil, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return
	}
	s.DebugPrint()
}

func insertion(c *Model.City) {
	s, err := local_search.Insertion(c, nil, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return
	}

	s.DebugPrint()
}

func threeOpt(c *Model.City) {
	s, err := local_search.ThreeOpt(c, nil, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return
	}
	s.DebugPrint()
}

func cheapest(c *Model.City) (*Model.Solution, bool) {
	s, err := builders.CheapestInsertion(c, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return nil, true
	}
	s.DebugPrint()
	return s, false
}

func nearest(c *Model.City) (*Model.Solution, error, bool) {
	s, err := builders.NearestNeighbor(c, nil)
	if err != nil {
		fmt.Println("\033[31m\033[1m", err, "\033[0m")
		return nil, nil, true
	}

	s.DebugPrint()
	return s, err, false
}
