package main

import "fmt"

func main() {
	batchPipelines()
	streamPipelines()
}

func streamPipelines() {
	multiply := func(value int, multiplier int) int {
		return value * multiplier
	}
	add := func(value int, additive int) int {
		return value + additive
	}
	ints := []int{1, 2, 3, 4}
	fmt.Println()
	for _, val := range ints {
		fmt.Printf("%d ", add(multiply(val, 2), 1))
	}
}

func batchPipelines() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}
	fmt.Println(multiply([]int{1, 2, 3, 4}, 2))

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}
	fmt.Println(add([]int{2, 4, 6, 8}, 1))

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Printf("%d ", v)
	}
}
