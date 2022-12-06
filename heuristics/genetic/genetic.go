package genetic

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"tcp_drone/heuristics/builders"
	local_search_drone "tcp_drone/heuristics/local_search_drone"
	Model "tcp_drone/model"
)

type GeneticSettings struct {
	// number of generations
	Generations int
	// mutation rate
	MutationRate float64
	// crossover rate
	CrossoverRate float64
	// elitism
	Elitism bool
	// elitism size
	ElitismSize int
	// selection type
	SelectionType string
	// crossover type
	CrossoverType string
	// mutation type
	MutationType string
	// population size
	PopulationSize int
}

// create model to genetic algorithm
type GeneticModel struct {
	//Settings
	Settings GeneticSettings
	// Population
	Population []Model.Solution
	// Best solution
	BestSolution Model.Solution
	// Selection
	Selection SelectionFunc
	// Crossover
	Crossover CrossOverFunc
	// Mutation
	Mutation MutationFunc
	//Stop condition
	StopCondition func(generation int, population []Model.Solution) bool
}

func GeneticAlgorithm(city *Model.City, solution *Model.Solution, builderHeuristic builders.BuilderHeuristic) (*Model.Solution, error) {
	if solution == nil {
		solution = &Model.Solution{}
		e := solution.Init(*city)
		if e != nil {
			return nil, e
		}
	}
	defer solution.Timer()()
	geneticAlgorithm(city, solution)
	return solution, nil
}

func geneticAlgorithm(city *Model.City, solution *Model.Solution) {
	if solution.RouteTime == 0 {
		_, _ = local_search_drone.VndDrone(city, solution, nil)
	}

	elitismSize := int(float64(len(city.Nodes)) * 0.1)

	// init genetic model
	geneticModel := GeneticModel{
		Settings: GeneticSettings{
			Generations:    100000,
			MutationRate:   0.5,
			CrossoverRate:  1,
			Elitism:        true,
			ElitismSize:    elitismSize,
			SelectionType:  "tournament",
			CrossoverType:  "order",
			MutationType:   "swap",
			PopulationSize: 100,
		},
		Population: make([]Model.Solution, 0),
	}

	// init population
	//geneticModel.Population = append(geneticModel.Population, *solution)
	for i := 0; i < geneticModel.Settings.PopulationSize; i++ {
		_solution := Model.Solution{GlobalSearchAlgorithm: "Genetic"}
		_solution.Init(*city)
		_solution.Route = generateRandomRoute(city)
		_solution.Fitness()
		geneticModel.Population = append(geneticModel.Population, _solution)
	}

	//printPopulation(geneticModel)

	// init selection
	switch geneticModel.Settings.SelectionType {
	case "tournament":
		geneticModel.Selection = tournamentSelection
		//case "roulette":
		//	geneticModel.Selection = rouletteSelection
		//case "rank":
		//	geneticModel.Selection = rankSelection
	}

	// init crossover
	switch geneticModel.Settings.CrossoverType {
	case "order":
		geneticModel.Crossover = orderCrossover
		//case "pmx":
		//	geneticModel.Crossover = pmxCrossover
		//case "ox":
		//	geneticModel.Crossover = ox1Crossover
	}

	// init mutation
	switch geneticModel.Settings.MutationType {
	case "swap":
		geneticModel.Mutation = swapMutation
		//case "insert":
		//	geneticModel.Mutation = insertMutation
		//case "scramble":
		//	geneticModel.Mutation = scrambleMutation
	}

	// init stop condition
	geneticModel.StopCondition = func(generation int, population []Model.Solution) bool {
		return generation >= geneticModel.Settings.Generations
	}

	// run genetic algorithm
	*solution = geneticModel.Run()
}

func (geneticModel *GeneticModel) Run() Model.Solution {
	geneticModel.Population = sortByFitness(geneticModel.Population)
	// init best solution
	geneticModel.BestSolution = Model.Solution{Route: make([]int, len(geneticModel.Population[0].Route)), City: geneticModel.Population[0].City}
	copy(geneticModel.BestSolution.Route, geneticModel.Population[0].Route)
	geneticModel.BestSolution.RouteTime = geneticModel.Population[0].RouteTime

	fmt.Println("Start Best solution: ", geneticModel.BestSolution.RouteTime)

	// init generation
	generation := 0

	// run genetic algorithm
	for !geneticModel.StopCondition(generation, geneticModel.Population) {

		// init new population
		newPopulation := make([]Model.Solution, 0)

		// init elitism
		if geneticModel.Settings.Elitism {
			elitism := make([]Model.Solution, 0)
			for i := 0; i < geneticModel.Settings.ElitismSize; i++ {
				elitism = append(elitism, geneticModel.Population[i])
			}
			newPopulation = append(newPopulation, elitism...)
		}

		wg := sync.WaitGroup{}
		iterations := (geneticModel.Settings.PopulationSize - geneticModel.Settings.ElitismSize) / 2
		wg.Add(iterations)
		for i := 0; i < iterations; i++ {
			go geneticModel.testeThread(&newPopulation, &wg)
		}
		wg.Wait()

		// update population
		geneticModel.Population = newPopulation

		for i := 0; i < len(geneticModel.Population); i++ {
			geneticModel.Population[i].Fitness()
		}

		// update best solution
		geneticModel.Population = sortByFitness(geneticModel.Population)
		if geneticModel.Population[0].RouteTime < geneticModel.BestSolution.RouteTime {
			geneticModel.BestSolution = Model.Solution{Route: make([]int, len(geneticModel.Population[0].Route)), City: geneticModel.Population[0].City}
			copy(geneticModel.BestSolution.Route, geneticModel.Population[0].Route)
			geneticModel.BestSolution.RouteTime = geneticModel.Population[0].RouteTime
		}

		// update generation
		generation++
	}

	return geneticModel.BestSolution
}

func (geneticModel *GeneticModel) testeThread(newPopulation *[]Model.Solution, wg *sync.WaitGroup) {

	// selection
	parent1 := geneticModel.Selection(geneticModel.Population)
	parent2 := geneticModel.Selection(geneticModel.Population)

	// crossover
	if rand.Float64() < geneticModel.Settings.CrossoverRate {
		parent1, parent2 = geneticModel.Crossover(parent1, parent2)
	}

	// mutation
	if rand.Float64() < geneticModel.Settings.MutationRate {
		geneticModel.Mutation(&parent1)
		geneticModel.Mutation(&parent2)
	}

	// add to new population
	*newPopulation = append(*newPopulation, parent1, parent2)
	wg.Done()
}

func sortByFitness(population []Model.Solution) []Model.Solution {

	sort.Slice(population, func(i, j int) bool {
		return population[i].RouteTime < population[j].RouteTime
	})

	return population
}

// generate random route solutions, first and last city is 0
func generateRandomRoute(city *Model.City) []int {
	route := make([]int, 0)
	for i := 0; i < len(city.Nodes); i++ {
		route = append(route, i)
	}
	rand.Shuffle(len(route), func(i, j int) {
		if i != j && i != 0 && j != 0 {
			route[i], route[j] = route[j], route[i]
		}
	})
	route = append(route, 0)
	return route
}

func printPopulation(geneticModel GeneticModel) {
	for i := 0; i < len(geneticModel.Population); i++ {
		fmt.Println(geneticModel.Population[i].Route, geneticModel.Population[i].RouteTime)
	}
}
