package genetic

import (
	"math/rand"
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
	Mutation func(solution *Model.Solution)
	//Stop condition
	StopCondition func(generation int, population []Model.Solution) bool
}

func GeneticAlgorithm(city *Model.City, solution *Model.Solution) (*Model.Solution, error) {
	defer solution.Timer()()
	geneticAlgorithm(city, solution)
	return solution, nil
}

func geneticAlgorithm(city *Model.City, solution *Model.Solution) {
	// init genetic model
	geneticModel := GeneticModel{
		Settings: GeneticSettings{
			Generations:    100,
			MutationRate:   0.1,
			CrossoverRate:  0.7,
			Elitism:        true,
			ElitismSize:    1,
			SelectionType:  "tournament",
			CrossoverType:  "order",
			MutationType:   "swap",
			PopulationSize: 100,
		},
		Population: make([]Model.Solution, 0),
	}

	// init population
	for i := 0; i < geneticModel.Settings.PopulationSize; i++ {
		solution := Model.Solution{BuilderAlgorithm: "Genetic"}
		solution.Init(*city)
		geneticModel.Population = append(geneticModel.Population, solution)
	}

	// init selection
	switch geneticModel.Settings.SelectionType {
	case "tournament":
		geneticModel.Selection = tournamentSelection
	case "roulette":
		geneticModel.Selection = rouletteSelection
	case "rank":
		geneticModel.Selection = rankSelection
	}

	// init crossover
	switch geneticModel.Settings.CrossoverType {
	case "order":
		geneticModel.Crossover = orderCrossover
	case "pmx":
		geneticModel.Crossover = pmxCrossover
	case "ox":
		geneticModel.Crossover = ox1Crossover
	}

	// init mutation
	switch geneticModel.Settings.MutationType {
	case "swap":
		geneticModel.Mutation = swapMutation
	case "insert":
		geneticModel.Mutation = insertMutation
	case "scramble":
		geneticModel.Mutation = scrambleMutation
	}

	// init stop condition
	geneticModel.StopCondition = func(generation int, population []Model.Solution) bool {
		return generation >= geneticModel.Settings.Generations
	}

	// run genetic algorithm
	geneticModel.Run(city)

}

func (geneticModel *GeneticModel) Run(city *Model.City) {
	// init best solution
	geneticModel.BestSolution = geneticModel.Population[0]

	// init generation
	generation := 0

	// run genetic algorithm
	for !geneticModel.StopCondition(generation, geneticModel.Population) {
		// init new population
		newPopulation := make([]Model.Solution, 0)

		// init fitness
		fitness := make([]float64, 0)

		// calculate fitness
		for _, solution := range geneticModel.Population {
			fitness = append(fitness, 1/solution.RouteTime)
		}

		// init elitism
		if geneticModel.Settings.Elitism {
			elitism := make([]Model.Solution, 0)
			for i := 0; i < geneticModel.Settings.ElitismSize; i++ {
				elitism = append(elitism, geneticModel.BestSolution)
			}
			newPopulation = append(newPopulation, elitism...)
		}

		// selection
		selected := geneticModel.Selection(geneticModel.Population, fitness)

		// crossover
		for i := 0; i < len(selected); i += 2 {
			parent1 := geneticModel.Population[selected[i]]
			parent2 := geneticModel.Population[selected[i+1]]
			child1, child2 := geneticModel.Crossover(parent1, parent2)
			newPopulation = append(newPopulation, Model.Solution{Route: child1}, Model.Solution{Route: child2})
		}

		// mutation
		for i := 0; i < len(newPopulation); i++ {
			if rand.Float64() < geneticModel.Settings.MutationRate {
				geneticModel.Mutation(&newPopulation[i])
			}
		}

		// calculate fitness
		for i := 0; i < len(newPopulation); i++ {
			newPopulation[i].Calculate()
		}

		// update population
		geneticModel.Population = newPopulation

		// update best solution
		for _, solution := range geneticModel.Population {
			if solution.RouteTime < geneticModel.BestSolution.RouteTime {
				geneticModel.BestSolution = solution
			}
		}

		// update generation
		generation++
	}
}
