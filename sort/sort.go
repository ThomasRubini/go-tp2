package main

import (
	"fmt"
	"math"
	"sort"
)

func splitTab(tab []int, N int) [][]int {
	parts := make([][]int, N)
	maxPartSize := int(math.Ceil(float64(len(tab)) / float64(N)))

	for i := 0; i < N; i++ {
		minIndex := i * maxPartSize
		maxIndex := int(math.Min(float64((i+1)*maxPartSize), float64(len(tab))))
		parts[i] = tab[minIndex:maxIndex]
	}

	return parts
}

func merge(tab []int, output chan []int) {
	sort.Ints(tab)
	output <- tab
}

func mergeSort(tab []int, N int) []int {
	parts := splitTab(tab, N)

	// Send work
	outputs := make([]chan []int, N)
	for i := 0; i < N; i++ {
		outputs[i] = make(chan []int)
		go merge(parts[i], outputs[i])
	}

	// Receive work
	sortedParts := make([][]int, N)
	for i := 0; i < N; i++ {
		sortedParts[i] = <-outputs[i]
	}

	// Sort
	sortedTab := make([]int, 0)
	for {
		// Get part with smallest leading number
		var candidateIndex *int
		for i := 0; i < N; i++ {
			if len(sortedParts[i]) > 0 && (candidateIndex == nil || sortedParts[i][0] < sortedParts[*candidateIndex][0]) {
				candidateIndex = &i
			}
		}

		// No candidate, everything is in the sorted tab
		if candidateIndex == nil {
			break
		}

		// Candidate, add it to the sorted tab and remove it from the part
		sortedTab = append(sortedTab, sortedParts[*candidateIndex][0])
		sortedParts[*candidateIndex] = sortedParts[*candidateIndex][1:]
	}

	return sortedTab
}

func main() {
	tab := []int{10, 5, 7, 2, 9, 1, 0, 7, 23}
	result := mergeSort(tab, 3)
	fmt.Println(result)
}
