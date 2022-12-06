package model

import (
	"math"
	array "tcp_drone/utils"
)

type SubRoute struct {
	OnlyDriver bool
	Route      []int
}

type SubSolution struct {
	Customer  int
	IsDrone   bool
	MaxSaving float64
	Savings   float64
	Departure int
	Arrival   int
}

// CalcSavings Calculate the savings of the solution change the route to drone
func (s *Solution) CalcSavings(customerRouteIndex int, subRoutes []SubRoute) float64 {
	var i = array.Abs(s.Route[customerRouteIndex-1])
	var j = array.Abs(s.Route[customerRouteIndex])
	var k = array.Abs(s.Route[customerRouteIndex+1])

	var t = s.City.GetDriverTimeDistance
	var td = s.City.GetDroneTimeDistance
	savings := t(i, j) + t(j, k) - t(i, j)
	for _, subRoute := range subRoutes {
		// printf warning bold and yellow
		if array.Contains(subRoute.Route, j) && !subRoute.OnlyDriver {
			a := subRoute.Route[0]
			b := subRoute.Route[len(subRoute.Route)-1]
			droneNode := subRoute.Route[1]
			_, _, driverTime, _, _ := CalculateDroneFlightTime(s.City, s.simulateSubRouteWithoutItem(subRoute, j), 1)
			savings = math.Min(savings, driverTime+td(a, droneNode)+td(droneNode, b)+s.City.DroneSettings.ReceiptTime)
			break
		}
	}
	return savings
}

// GetEligibleDroneCustomers get elegibles drone customers
func (c *City) GetEligibleDroneCustomers(excludeNodes *[]int) []Customer {
	if excludeNodes == nil {
		excludeNodes = &[]int{}
	}
	var eligibleCustomers []Customer
	for _, n := range c.Nodes {
		if array.Contains(*excludeNodes, n.Node) || !n.IsAllowedDrone || n.Node == 0 {
			continue
		}
		eligibleCustomers = append(eligibleCustomers, n)
	}
	return eligibleCustomers
}

// GetSubRoutes get the subroute that contains drone and driver
func (s *Solution) GetSubRoutes() []SubRoute {
	//	create and 2d array
	var subRoutes []SubRoute
	var actualSubRoute = SubRoute{OnlyDriver: true, Route: []int{}}
	for i, e := range s.Route {
		if e >= 0 {
			if !actualSubRoute.OnlyDriver {
				actualSubRoute.Route = append(actualSubRoute.Route, e)
				subRoutes = append(subRoutes, actualSubRoute)
				actualSubRoute = SubRoute{OnlyDriver: true, Route: []int{}}
			}
			actualSubRoute.Route = append(actualSubRoute.Route, e)
		} else {
			if actualSubRoute.OnlyDriver {
				subRoutes = append(subRoutes, actualSubRoute)
				actualSubRoute = SubRoute{OnlyDriver: false, Route: []int{}}
				actualSubRoute.Route = append(actualSubRoute.Route, s.Route[i-1])
			}
			actualSubRoute.Route = append(actualSubRoute.Route, e)

		}
	}
	subRoutes = append(subRoutes, actualSubRoute)
	return subRoutes
}

// AppendSubRoute append subroute to the solution
func AppendSubRoute(subRoutes []SubRoute) []int {
	var newRoute []int
	for i, subRoute := range subRoutes {
		if i == 0 {
			newRoute = append(newRoute, subRoute.Route...)
		} else {
			newRoute = append(newRoute, subRoute.Route[1:]...)
		}
	}

	return newRoute
}

// CalcCostDriver Calculate the cost of the solution with only driver
func (s *Solution) CalcCostDriver(subRoute SubRoute, customer int, subSolution *SubSolution) float64 {
	t := s.City.GetDriverTimeDistance
	j := customer
	// copy the route
	routeWithoutDroneNode := append([]int(nil), subRoute.Route...)
	routeWithoutDroneNode = append(routeWithoutDroneNode[:1], routeWithoutDroneNode[2:]...)
	for l, e := range routeWithoutDroneNode {
		if l+1 < len(routeWithoutDroneNode) {
			i := array.Abs(e)
			k := array.Abs(routeWithoutDroneNode[l+1])
			cost := t(i, j) + t(j, k) - t(i, k)
			if cost < subSolution.Savings {
				_, _, driveTime, _, _ := CalculateDroneFlightTime(s.City, s.simulateSubRouteWithItem(subRoute.Route, j, i), 1)
				if driveTime+cost <= s.City.DroneSettings.BatteryTime && subSolution.Savings-cost > subSolution.MaxSaving {
					subSolution.Departure = i
					subSolution.Customer = j
					subSolution.IsDrone = false
					subSolution.MaxSaving = subSolution.Savings - cost
				}
			}
		}
	}
	// print max saving blue
	return subSolution.MaxSaving
}

// SimulateSubRouteWithItem simulate subroute without item
func (s *Solution) simulateSubRouteWithItem(route []int, customer int, afterCustomer int) []int {
	var newRoute []int
	if afterCustomer == route[0] {
		afterCustomer = route[1]
	}
	for _, e := range route {
		newRoute = append(newRoute, e)
		if e == afterCustomer {
			newRoute = append(newRoute, -customer)
		}
	}
	return newRoute
}

func (s *Solution) simulateSubRouteWithoutItem(subRoute SubRoute, j int) []int {
	var newRoute []int
	if subRoute.Route[0] == j || subRoute.Route[len(subRoute.Route)-1] == j {
		newRoute = append(newRoute, subRoute.Route...)
		index, _ := s.FindCustomerIndex(j)
		if subRoute.Route[0] == j {
			newItem := s.Route[index-1]
			newRoute = array.InsertAt(newRoute, 0, newItem)
		} else {
			newItem := s.Route[index+1]
			newRoute = append(newRoute, newItem)
		}
	} else {
		newRoute = array.Remove(subRoute.Route, j)
	}
	return newRoute
}

// CalcCostDrone Calculate the cost of the solution with only drone
func (s *Solution) CalcCostDrone(subRoute SubRoute, customer int, subSolution *SubSolution) float64 {
	//td := s.City.GetDroneTimeDistance
	sl := s.City.DroneSettings.LaunchTime
	sr := s.City.DroneSettings.ReceiptTime
	//batteryTime := s.City.DroneSettings.BatteryTime
	j := customer
	for l, i := range subRoute.Route {
		if i == j {
			continue
		}
		for h := l + 1; h < len(subRoute.Route); h++ {
			k := subRoute.Route[h]
			if k == j {
				continue
			}
			// create a new route from subroute that init on l until h
			var newRoute []int
			for u := l; u < h+1; u++ {
				if subRoute.Route[u] == j {
					continue
				}
				newRoute = append(newRoute, subRoute.Route[u])
			}
			newRoute = array.InsertAt(newRoute, 1, -j)
			_, isValid, driveTime, droneTime, _ := CalculateDroneFlightTime(s.City, newRoute, 1)
			if isValid {
				_cost := math.Max(driveTime, droneTime) + sl + sr - driveTime
				cost := math.Max(0, _cost)
				if subSolution.Savings-cost > subSolution.MaxSaving {
					subSolution.Customer = -j
					subSolution.Departure = i
					subSolution.Arrival = k
					subSolution.IsDrone = true
					subSolution.MaxSaving = subSolution.Savings - cost
				}
			}
		}
	}
	return subSolution.MaxSaving
}

// SwapCustomer swap customer
func (s *Solution) SwapCustomer(customer int, preCustomer int) {
	//TODO: Optimize
	//s.Route = array.Remove(s.Route, array.Abs(customer))
	//index, _ := s.FindCustomerIndex(preCustomer)
	//s.Route = array.InsertAt(s.Route, index+1, customer)

	c := array.Abs(customer)
	exclude := false
	findIndex := false
	var index int
	for i, e := range s.Route {
		if array.Abs(e) == c {
			s.Route = append(s.Route[:i], s.Route[i+1:]...)
			exclude = true

			if s.Route[i] == preCustomer {
				index = i
				findIndex = true
			}
		} else if e == preCustomer {
			index = i
			findIndex = true
		}
		if exclude && findIndex {
			break
		}
	}

	if s.Route[index] == preCustomer {
		s.Route = array.InsertAt(s.Route, index+1, customer)
	} else if s.Route[index-1] == -preCustomer {
		s.Route = array.InsertAt(s.Route, index, customer)
	} else if s.Route[index+1] == -preCustomer {
		s.Route = array.InsertAt(s.Route, index+2, customer)
	} else {
		panic("error")
	}
}
