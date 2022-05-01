package main

import (
	"fmt"
)

func main() {
	var pebbleWeight float64 = 0.1
	var rockWeight float64 = 1.2
	var boulderWeight float64 = 502.4
	var totalWeight float64 = pebbleWeight + rockWeight + boulderWeight
	fmt.Println(totalWeight)
}