package model

import (
	"fmt"
	"math/rand"
)

func Test() {
	//Testing city
	println("Testing city")
	c := &City{}
	c.Init("tsplib_generated_instances/berlin52/", nil, nil)
	// pirint index, x, y, isAllowedDrone for 17
	index := 0
	c.Nodes[index].Print()
	fmt.Println("NodesLen: ", c.NodesLen)
	fmt.Println("DriverDistance[0][1]: ", c.DriverTimeDistance[0][51])
	fmt.Println("DroneDistance[0][1]: ", c.DroneTimeDistance[0][1])

	//Testing solution
	println("Testing solution")
	s := &Solution{}
	s.Init(*c)
	// populate route with nodesLen+2
	for i := 0; i <= c.NodesLen; i++ {
		if i == 0 || i == c.NodesLen {
			s.Route[i] = 0
		} else {
			s.Route[i] = i
		}
	}
	fmt.Println("Route: ", s.Route)
	// shuffle route, except first and last
	for i := 1; i < c.NodesLen-1; i++ {
		j := rand.Intn(c.NodesLen-1) + 1
		s.Route[i], s.Route[j] = s.Route[j], s.Route[i]
	}
	// calculate time
	s.Calculate()
	s.DebugPrint()

	s.Route = []int{0, 21, -31, -48, -35, -34, -33, 38, -39, -37, -4, -14, -5, 23, 47, -36, -45, -43, -15, 49, 19, 22, 30, 17, 2, -44, -18, -40, -7, 8, 9, 42, 3, 24, 11, 27, 26, 25, 46, 13, 12, -51, -10, -50, 32, 28, 29, 20, 16, 41, 6, 1, 0}
	s.Calculate()
	s.DebugPrint()

	//s.Route = []int{0, 21, 31, 48, 35, 34, 33, 38, 39, 37, 4, 14, 5, 23, 47, 36, 45, 43, 15, 49, 19, 22, 30, 17, 2, 44, 18, 40, 7, 8, 9, 32, 42, -3, 24, 11, 27, 26, 25, 46, 13, 12, 51, 10, 50, 28, 29, 6, 1, 41, 16, 20, 0}
	//
	//subRoutes := s.GetSubRoutes()
	//fmt.Println("SubRoutes: ", subRoutes)
	//
	//s.Route[1] = -s.Route[1]
	//s.Route[2] = -s.Route[2]
	//s.Route[4] = -s.Route[4]
	//s.Route[len(s.Route)-2] = -s.Route[len(s.Route)-2]
	//s.Route[len(s.Route)-3] = -s.Route[len(s.Route)-3]
	//fmt.Println("Route: ", s.Route)
	//subRoutes = s.GetSubRoutes()
	//fmt.Println("SubRoutes: ", subRoutes)
	//
	//s.Calculate()
	//s.DebugPrint()

}
