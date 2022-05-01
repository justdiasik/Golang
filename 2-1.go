package main

import (
	"fmt"
	"strings"
)

func scoreSummary(name string, score1, score2, score3 float64) float64 {
	return (score1 + score2 + score3) / 3.0
}

func main() {
	var names = []string{"Jermaine", "Jessie", "Lamar"}
	var scores = [][]float64{
		{95.40, 82.30, 74.60},
		{79.30, 99.10, 82.50},
		{82.20, 95.40, 77.60},
	}
	headers := []string{"Name", "Score 1", "Score 2", "Score 3", "Average"}
	headerFormat := "%10s|%10s|%10s|%10s|%10s\n"
	rowFormat := "%10s|%10.02f|%10.02f|%10.02f|%10.02f\n"
	headerRow := fmt.Sprintf(headerFormat, headers[0], headers[1], headers[2], headers[3], headers[4])
	fmt.Print(headerRow)
	fmt.Println(strings.Repeat("-", len(headerRow)))
	for i := range names {
		scoreSummary := scoreSummary(names[i], scores[i][0], scores[i][1], scores[i][2])
		fmt.Printf(rowFormat, names[i], scores[i][0], scores[i][1], scores[i][2], scoreSummary)
	}
}
