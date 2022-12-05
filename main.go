package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strconv"
	"sync"
	"tcp_drone/heuristics"
	"tcp_drone/heuristics/builders_drone"
	"tcp_drone/heuristics/global_search"
	local_search "tcp_drone/heuristics/local_search"
	local_search_drone "tcp_drone/heuristics/local_search_drone"
	"tcp_drone/model"
)

type Instance struct {
	Name      string
	Path      string
	Reference float64
	Solutions []model.Solution
	City      model.City
}

func main() {
	//model.Test()
	heuristics.Test()

	//instances, header := readInstances()
	//appendHeader := func(s string) {
	//	header = append(header, s)
	//}
	//processInstances(instances, appendHeader)
	//renderTable(instances, header)
}

func readInstances() ([]Instance, []string) {
	path := "tsplib_generated_instances/"
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	instances := make([]Instance, 0)
	references := []float64{1000, 210.03, 3456.80, 178.16, 461.83, 1000, 13.45, 16.35, 587.80, 764.42, 870.65, 763.15, 835.43, 658.38, 606.45, 651.31, 378.25, 1204.42, 1653.80, 2642.00, 1666.25, 2114.04, 71.40, 37.15, 240.46, 20.50}
	count := 0
	for _, f := range files {
		if f.IsDir() {
			path := path + f.Name() + "/"
			instances = append(instances, Instance{Name: f.Name(), Path: path, City: model.City{}, Reference: references[count]})
			instances[len(instances)-1].City.Init(path, nil, nil)
			count++
		}
	}
	header := []string{"Index", "Name", "Nodes", "Ref"}
	return instances, header
}

// create a print

func processInstances(instances []Instance, header func(s string)) {
	header("VND")
	header("Drone")
	header("VND-Drone")
	//header("VND-Drone")
	header("SA-Drone")
	var wg sync.WaitGroup
	wg.Add(len(instances))

	for i := range instances {
		go engine(instances, i, &wg)
	}
	wg.Wait()
}

func engine(instances []Instance, i int, group *sync.WaitGroup) {
	defer group.Done()
	fmt.Println("Processing instance: " + instances[i].Name)
	var s *model.Solution

	//s, _ = builders.NearestNeighbor(&instances[i].City, nil)
	//instances[i].Solutions = append(instances[i].Solutions, *s)
	//
	sVnd, _ := local_search.Vnd(&instances[i].City, nil, nil)
	instances[i].Solutions = append(instances[i].Solutions, *sVnd)

	sDrone, _ := builders_drone.DroneInitialSolution(&instances[i].City, sVnd, nil)
	instances[i].Solutions = append(instances[i].Solutions, *sDrone)

	sVndD, _ := local_search_drone.VndDrone(&instances[i].City, sDrone, nil)
	instances[i].Solutions = append(instances[i].Solutions, *sVndD)

	s, _ = global_search.SimulatedAnnealing(&instances[i].City, sVndD, nil)
	instances[i].Solutions = append(instances[i].Solutions, *s)

	fmt.Println("Finished instance: " + instances[i].Name + " with solution: " + strconv.FormatFloat(s.RouteTime, 'f', 2, 64))
}

func renderTable(instances []Instance, header []string) {
	data := make([][]string, 0)

	for i := range instances {
		arr := make([]string, 0)
		arr = append(arr, strconv.Itoa(i))
		arr = append(arr, instances[i].Name)
		arr = append(arr, strconv.Itoa(instances[i].City.NodesLen))
		arr = append(arr, strconv.FormatFloat(instances[i].Reference, 'f', 2, 64))
		for j := range instances[i].Solutions {
			// if route time better than reference print green and bold, yellow and bold if is 20%, red and bold if much worse
			st := strconv.FormatFloat(instances[i].Solutions[j].RouteTime, 'f', 2, 64)
			set := instances[i].Solutions[j].ExecutionTime.String()
			if instances[i].Solutions[j].RouteTime <= instances[i].Reference {
				arr = append(arr, "\033[32m\033[1m"+st+"\033[0m"+"/"+set)
			} else {
				arr = append(arr, "\033[31m\033[1m"+st+"\033[0m"+"/"+set)
			}
		}
		data = append(data, arr)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetCenterSeparator("╬")
	table.SetColumnSeparator("║")
	table.SetRowSeparator("═")
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}

	table.Render()
}
