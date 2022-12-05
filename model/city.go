package model

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

import array "tcp_drone/utils"

// DroneSettings set default values
func (c *City) defaultDroneSettings() {
	c.DroneSettings.LaunchTime = 0.01667
	c.DroneSettings.ReceiptTime = 0.01667
	c.DroneSettings.BatteryTime = 0.67
}

type Customer struct {
	Node           int
	X              int64
	Y              int64
	IsAllowedDrone bool
}

type DroneSettings struct {
	LaunchTime  float64
	ReceiptTime float64
	BatteryTime float64
	Resistance  float64
	Penalty     float64
}

type City struct {
	DriverTimeDistance [][]float64
	DroneTimeDistance  [][]float64
	Nodes              []Customer
	NodesLen           int
	DroneSettings      DroneSettings
}

/// Main functions

// Init initialize city from csv files
func (c *City) Init(directoryPath string, droneSettings *DroneSettings, removeLast *bool) {
	//check if directory exists
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		//print red
		fmt.Printf("\033[31m%s\033[0m\n", "Directory does not exist")
		panic(err)
	}
	//create city
	c.loadNodesFromCSV(directoryPath+"nodes.csv", removeLast)
	c.loadDriverDistanceFromCSV(directoryPath + "tau.csv")
	c.loadDroneDistanceFromCSV(directoryPath + "tauprime.csv")

	if droneSettings == nil {
		c.defaultDroneSettings()
		//fmt.Printf("\033[31m%s\033[0m\n", "Drone settings not provided, using default values")
		//fmt.Printf("\033[31m%f\033[0m\n", c.DroneSettings.LaunchTime)
	} else {
		c.DroneSettings = *droneSettings
	}
}

// calcDroneTimeDistance calculate drone time distance
func (c *City) calcDroneTimeDistance(i, j, k int) float64 {
	cacheDroneTime := c.DroneSettings.LaunchTime + c.DroneTimeDistance[i][j] + c.DroneTimeDistance[j][k] + c.DroneSettings.ReceiptTime
	cacheDriverTime := c.DriverTimeDistance[j][k]
	if c.DroneSettings.BatteryTime < cacheDroneTime {
		return math.Inf(1)
	} else {
		return -math.Max(cacheDroneTime, cacheDriverTime)
	}
}

func (c *City) loadNodesFromCSV(fileName string, removeLast *bool) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := csv.NewReader(file)
	for {
		record, err := scanner.Read()
		if err != nil {
			break
		}
		var node Customer
		// print record
		node.Node, _ = strconv.Atoi(record[0])
		node.X, _ = strconv.ParseInt(strings.TrimSpace(record[1]), 10, 64)
		node.Y, _ = strconv.ParseInt(strings.TrimSpace(record[2]), 10, 64)
		node.IsAllowedDrone = record[3] == "0" // 1 means that the package is too heavy for drone
		c.Nodes = append(c.Nodes, node)
	}
	c.NodesLen = len(c.Nodes)
}

func (c *City) loadDriverDistanceFromCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := csv.NewReader(file)
	for {
		record, err := scanner.Read()
		if err != nil {
			break
		}
		var row []float64
		for _, s := range record {
			f, _ := strconv.ParseFloat(s, 64)
			row = append(row, f)
		}
		c.DriverTimeDistance = append(c.DriverTimeDistance, row)
	}
}

func (c *City) loadDroneDistanceFromCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			println("Error closing file")
		}
	}(file)

	scanner := csv.NewReader(file)
	for {
		record, err := scanner.Read()
		if err != nil {
			break
		}
		var row []float64
		for _, s := range record {
			f, _ := strconv.ParseFloat(s, 64)
			row = append(row, f)
		}
		c.DroneTimeDistance = append(c.DroneTimeDistance, row)
	}
}

// GetDriverTimeDistance get driver time distance
func (c *City) GetDriverTimeDistance(i, j int) float64 {
	return c.DriverTimeDistance[array.Abs(i)][array.Abs(j)]
}

// GetDroneTimeDistance get drone time distance
func (c *City) GetDroneTimeDistance(i, j int) float64 {
	return c.DroneTimeDistance[array.Abs(i)][array.Abs(j)]
}

//Auxiliary functions to Nearest Neighbour algorithm

// getNearestNode get nearest node
func (c *City) getNearestNode(node int, targetMeasure func(int, int) float64, excludeNodes *[]int) Customer {
	if excludeNodes == nil {
		excludeNodes = &[]int{}
	}
	var nearestNode Customer
	var minDistance float64
	for _, n := range c.Nodes {
		if n.Node == node || array.Contains(*excludeNodes, n.Node) {
			continue
		}
		distance := targetMeasure(node, n.Node)
		if minDistance == 0 || distance < minDistance {
			minDistance = distance
			nearestNode = n
		}
	}
	return nearestNode
}

// GetNearestNodeDriver get nearest node
func (c *City) GetNearestNodeDriver(index int, excludeNodes *[]int) Customer {
	return c.getNearestNode(index, c.GetDriverTimeDistance, excludeNodes)
}

// GetNearestNodeDrone get nearest node
func (c *City) GetNearestNodeDrone(index int, excludeNodes *[]int) Customer {
	return c.getNearestNode(index, c.GetDroneTimeDistance, excludeNodes)
}

// GetCheapestNodeDriver get nearest node
func (c *City) GetCheapestNodeDriver(index int, excludeNodes *[]int) Customer {
	return c.getNearestNode(index, func(i, j int) float64 {
		return c.GetDriverTimeDistance(i, j) * c.DroneSettings.Resistance
	}, excludeNodes)
}

// Debug functions

// Print Customer print
func (c *Customer) Print() {
	fmt.Println("Nodes -> index:", c.Node, "x:", c.X, "y:", c.Y, "isAllowedDrone:", c.IsAllowedDrone)
}
