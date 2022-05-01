package main

import (
	"fmt"
)

func perimeter(width, height float64) float64 {
	return 2 * (width + height)
}

func main() {
	var widths = []float64{8.2, 5, 6.2}
	var heights = []float64{10, 5.2, 4.5}
	var total float64
	for i := range widths {
		total += perimeter(widths[i], heights[i])
	}
	fmt.Printf("%.02f", total)
}
