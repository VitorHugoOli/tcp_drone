package genetic

import (
	"math/rand"
	Model "tcp_drone/model"
	"tcp_drone/utils"
)

type CrossOverFunc func(solution1, solution2 Model.Solution, params ...interface{}) (Model.Solution, Model.Solution)

// orderCrossover performs order crossover on two solutions and returns two new solutions
func orderCrossover(solution1, solution2 Model.Solution, params ...interface{}) (Model.Solution, Model.Solution) {
	var dice int
	for dice%2 != 0 || dice > 2 || dice > len(solution1.City.Nodes)/2 {
		dice = rand.Intn(len(solution1.City.Nodes))
	}

	var cities []int
	for len(cities) < dice {
		city := rand.Intn(len(solution1.City.Nodes))
		if !utils.Contains(cities, city) && city != 0 {
			cities = append(cities, city)
		}
	}

	//create new solutions
	var (
		newSolution1 = Model.Solution{Route: make([]int, len(solution1.Route)), City: solution1.City}
		newSolution2 = Model.Solution{Route: make([]int, len(solution2.Route)), City: solution2.City}
	)

	//copy route from parent solutions
	copy(newSolution1.Route, solution1.Route)
	copy(newSolution2.Route, solution2.Route)

	swapCities(&newSolution1, &newSolution2, cities)

	newSolution1.Fitness()
	newSolution2.Fitness()

	return newSolution1, newSolution2
}

func swapCities(solution1, solution2 *Model.Solution, cities []int) {
	posCitiesSolution1 := make([]int, len(cities))
	posCitiesSolution2 := make([]int, len(cities))

	for i, e := range cities {
		for i2, e2 := range solution1.Route {
			if e == utils.Abs(e2) {
				posCitiesSolution1[i] = i2
			}
		}

		for i2, e2 := range solution2.Route {
			if e == utils.Abs(e2) {
				posCitiesSolution2[i] = i2
			}
		}
	}

	for i := 0; i < len(cities); i = i + 2 {
		solution1.Route[posCitiesSolution1[i]] = solution2.Route[posCitiesSolution2[i+1]]
		solution2.Route[posCitiesSolution2[i+1]] = solution1.Route[posCitiesSolution1[i]]
		solution1.Route[posCitiesSolution1[i+1]] = solution2.Route[posCitiesSolution2[i]]
		solution2.Route[posCitiesSolution2[i]] = solution1.Route[posCitiesSolution1[i+1]]
	}
}
