package model

import (
	"errors"
	"fmt"
	"math"
	"tcp_drone/utils"
	"time"
)

type Solution struct {
	BuilderAlgorithm      string
	LocalSearchAlgorithm  string
	GlobalSearchAlgorithm string
	ExecutionTime         time.Duration
	City                  City
	// Route is a list of indexes of the nodes
	// the negative values means that customer will be served by a drone
	// the positive values means that customer will be served by a driver
	Route      []int
	RouteTime  float64
	Iterations int
}

func (s *Solution) Init(city City) error {
	s.City = city
	// check if the route is already initialized
	if s.City.Nodes == nil {
		return errors.New("route is not initialized")
	}
	s.Route = make([]int, city.NodesLen+1) //+2 for depot
	s.Route[0] = 0
	s.Route[city.NodesLen] = 0
	s.RouteTime = 0
	return nil
}

// Calculate Calculate the time of the solution
// considering the time of the driver and the drone
func (s *Solution) Calculate() {
	s.RouteTime = CalculateRouteTime(s.City, s.Route, nil)
}

func (s *Solution) FindCustomerIndex(customer int) (int, bool) {
	for i, c := range s.Route {
		if utils.Abs(c) == utils.Abs(customer) {
			if c < 0 {
				return i, true
			}
			return i, false
		}
	}
	return -1, false
}

// DebugPrint print the solution
func (s *Solution) DebugPrint() {
	// print solution in bold and blue
	fmt.Printf("\033[1;34m%s %s\033[0m\n", "Solution:", s.BuilderAlgorithm+" > "+s.LocalSearchAlgorithm+" > "+s.GlobalSearchAlgorithm)
	fmt.Printf("Execution time %s\n", s.ExecutionTime)
	fmt.Println("Route: ", s.Route)
	fmt.Println("RouteTime: ", s.RouteTime)
	println("Iterations: ", s.Iterations)
	s.DebugPrintRouteDetails()
}

// DebugPrintRouteDetails print the route details
func (s *Solution) DebugPrintRouteDetails() {
	// print route passing by the customers
	fmt.Println("Route details: ")
	buffer := ""
	var bufferNextLine []rune
	t := 0.0
	cost := 0.0
	isFlying := false
	flyingCost := 0.0
	driverCacheTime := 0.0
	nodeDeparture := -1
	nodeArrival := -1
	lenMaxDistance := 0
	lastDroneNode := -1
	lastDriverNode := -1
	costStr := ""
	for i, customer := range s.Route {
		if customer < 0 {
			buffer += fmt.Sprintf("\033[1;33m%d\033[0m", customer)
			// "  " append the next line
			bufferNextLine = append(bufferNextLine, []rune("--")...)
		} else {
			buffer += fmt.Sprintf("\033[1;34m%d\033[0m", customer)
			bufferNextLine = append(bufferNextLine, []rune(" ")...)
		}

		if i+1 != len(s.Route) {
			nextCustomer := s.Route[i+1]
			customer = int(math.Abs(float64(customer)))
			nextCustomer = int(math.Abs(float64(nextCustomer)))

			if s.Route[i+1] >= 0 {
				if isFlying {
					isFlying = false
					nodeArrival = nextCustomer

					cost = s.City.GetDroneTimeDistance(lastDroneNode, nextCustomer)
					driverCacheTime += s.City.GetDriverTimeDistance(lastDriverNode, nextCustomer)
					flyingCost += cost + s.City.DroneSettings.ReceiptTime

					buffer += fmt.Sprintf("\u001B[0;33m - %.2f + %.2f - \u001B[0m", cost, s.City.DroneSettings.ReceiptTime)
					costStr = fmt.Sprintf("%.2f + %.2f", cost, s.City.DroneSettings.ReceiptTime)
				} else {
					cost = s.City.GetDriverTimeDistance(customer, nextCustomer)
					buffer += fmt.Sprintf(" - %.2f - ", cost)
					costStr = fmt.Sprintf("%.2f", cost)
					t += cost
				}
			} else if s.Route[i+1] < 0 {
				if !isFlying {
					isFlying = true
					nodeDeparture = customer
					lastDriverNode = customer
					lastDroneNode = nextCustomer
					driverCacheTime = 0.0

					cost = s.City.GetDroneTimeDistance(customer, nextCustomer)
					flyingCost = cost + s.City.DroneSettings.LaunchTime

					buffer += fmt.Sprintf("\u001B[0;33m - %.2f + %.2f - \u001B[0m", cost, s.City.DroneSettings.LaunchTime)
					costStr = fmt.Sprintf("%.2f + %.2f", cost, s.City.DroneSettings.LaunchTime)
				} else if isFlying && s.Route[i-1] == nodeDeparture {
					lastDriverNode = nextCustomer
					driverCacheTime += s.City.GetDriverTimeDistance(nodeDeparture, nextCustomer)

					buffer += fmt.Sprintf("\u001B[0;33m - x - \u001B[0m")
					costStr = fmt.Sprintf("x")
				} else {
					lastDriverNode = nextCustomer

					cost = s.City.GetDriverTimeDistance(customer, nextCustomer)
					driverCacheTime += cost

					buffer += fmt.Sprintf("\u001B[0;33m - %.2f - \u001B[0m", cost)
					costStr = fmt.Sprintf("%.2f", cost)
				}
			}

			if customer == nodeDeparture || customer == nodeArrival {
				bufferNextLine[len(bufferNextLine)-1] = '|'
			}

			if nextCustomer == nodeArrival {
				max := math.Max(flyingCost, driverCacheTime)
				maxString := fmt.Sprintf("%.2f", max)

				if max == flyingCost {
					maxString = "F:" + maxString
				} else {
					maxString = "D:" + maxString
				}

				// remove the length of the maxString from the bufferNextLine
				bufferNextLine = append(bufferNextLine, []rune(maxString)...)
				lenMaxDistance = len(maxString)
			}

			for j := 0; j < len(costStr)+6-lenMaxDistance; j++ {
				if isFlying || customer == nodeDeparture || nextCustomer == nodeArrival {
					bufferNextLine = append(bufferNextLine, []rune("\u001B[0;33m-\u001B[0m")...)
				} else {
					bufferNextLine = append(bufferNextLine, []rune(" ")...)
				}
			}
			lenMaxDistance = 0

		} else {
			if customer == nodeArrival {
				bufferNextLine[len(bufferNextLine)-1] = '|'
			}
		}

	}
	fmt.Println(buffer)
	fmt.Println(string(bufferNextLine))
}

func (s *Solution) Timer() func() {
	start := time.Now()
	return func() {
		s.ExecutionTime = time.Since(start)
		s.Calculate()
	}
}

// check if route have duplicate itens
func (s *Solution) CheckDuplicateItemsRoute() bool {
	visited := make(map[int]bool)
	for _, v := range s.Route {
		if v == 0 {
			continue
		}
		if visited[utils.Abs(v)] {
			return false
		}
		visited[utils.Abs(v)] = true
	}
	return true
}

// check if the route contains all nodes from s.City.Nodes
func (s *Solution) CheckAllItemsIndRoute() bool {
	visited := make(map[int]bool)
	for _, v := range s.Route {
		visited[utils.Abs(v)] = true
	}
	for _, v := range s.City.Nodes {
		if !visited[v.Node] {
			return false
		}
	}
	return true
}

func (s *Solution) CheckRoute() bool {
	if len(s.Route) == 0 {
		panic("route is empty")
	}
	if s.Route[0] != 0 || s.Route[len(s.Route)-1] != 0 {
		panic("route is not closed")
	}
	if !s.CheckAllItemsIndRoute() {
		panic("route does not contain all items")
	}
	if !s.CheckDuplicateItemsRoute() {
		panic("route contains duplicate items")
	}
	return true
}
