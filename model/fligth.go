package model

import (
	"math"
	"tcp_drone/utils"
)

func CalculateRouteTime(city City, route []int, isAllowedWrong *bool) float64 {
	t := 0.0
	driverLastPos := 0
	for i := 1; i < len(route)-1; i++ {
		customer := int(math.Abs(float64(route[i] * -1)))
		if route[i] < 0 {
			flightCost, isValid, _, _, jump := CalculateDroneFlightTime(city, route, i)
			if !isValid {
				return math.MaxFloat64
			}
			t += flightCost
			i = jump
			driverLastPos = route[i]
		} else {
			t += city.GetDriverTimeDistance(driverLastPos, customer)
			driverLastPos = route[i]
		}
	}
	return t
}

func CalculateRouteTimeGenetic(city City, route []int, isAllowedWrong *bool) float64 {
	t := 0.0
	driverLastPos := 0
	hasDrone := false
	for i := 1; i < len(route)-1; i++ {
		customer := int(math.Abs(float64(route[i] * -1)))
		if route[i] < 0 {
			hasDrone = true
			flightCost, isValid, _, _, jump := CalculateDroneFlightTime(city, route, i)
			if !isValid {
				return math.MaxFloat64
			}
			t += flightCost
			i = jump
			driverLastPos = route[i]
		} else {
			t += city.GetDriverTimeDistance(driverLastPos, customer)
			driverLastPos = route[i]
		}
	}

	//penalize if the solution does not use the drone
	if !hasDrone {
		t += 1000
	}
	return t
}

func CalculateDroneFlightTime(city City, route []int, startIndex int) (float64, bool, float64, float64, int) {
	customer := utils.Abs(route[startIndex])
	lr := city.DroneSettings.LaunchTime + city.DroneSettings.ReceiptTime
	var t = city.GetDriverTimeDistance
	var td = city.GetDroneTimeDistance
	cacheDroneTime := td(route[startIndex-1], customer)
	cacheDriveTime := t(route[startIndex-1], utils.Abs(route[startIndex+1]))
	nodeArrived := startIndex + 1
	for i := startIndex + 1; i < len(route) && route[i] < 0; i++ {
		node := utils.Abs(route[i])
		cacheDriveTime += t(node, route[i+1])
		nodeArrived = i + 1
	}
	cacheDroneTime += td(customer, route[nodeArrived])
	max := math.Max(cacheDroneTime, cacheDriveTime+lr)

	return max, max <= city.DroneSettings.BatteryTime, cacheDriveTime, cacheDroneTime, nodeArrived

}
