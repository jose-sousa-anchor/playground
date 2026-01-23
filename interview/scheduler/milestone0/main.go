package main

import "fmt"

// -----------------------------
// Warmup: processJob
// -----------------------------
// processJob takes a list of strings and returns a map where keys are string lengths
// and values are lists of strings with that length.
func processJob(items []string) map[int][]string {
	result := make(map[int][]string)
	for _, s := range items {
		l := len(s)
		result[l] = append(result[l], s)
	}
	return result
}

func main() {
	fmt.Println("=== Milestone 0: Warmup - processJob ===")
	ex1 := []string{"apple", "banana", "kiwi", "grape", "fig", "pear", "peach"}
	fmt.Println("Input:", ex1)
	fmt.Println("Grouped by length:", processJob(ex1))

	fmt.Println("\nEdge cases:")
	fmt.Println("Empty list:", processJob([]string{}))
	fmt.Println("With duplicates:", processJob([]string{"a", "aa", "a", "aaa"}))
}
