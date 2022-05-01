package main

import (
	"errors"
	"fmt"
)

func perimeter2(width, height float64) (float64, error) {
	if width < 0 {
		return 0, errors.New("a length of [length] is invalid")
	}
	if height < 0 {
		return 0, errors.New("a width of [width] is invalid")
	}
	return 2 * (width + height), nil
}

func main() {
	var widths = []float64{8.2, 5, 6.2, -1}
	var heights = []float64{10, 5.2, 4.5, 5.2}
	var total float64
	for i := range widths {
		result, err := perimeter2(widths[i], heights[i])
		if err != nil {
			panic(err)
		}
		total += result
	}
	fmt.Printf("%.02f", total)
}
