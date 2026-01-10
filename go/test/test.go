package main

import (
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
)

func printSolution(sol verifiers.Solution) {
	println("Solution:")
	fmt.Printf("Display:")
	for _, c := range sol.Display {
		fmt.Printf("%c", c)
	}
	fmt.Printf("\nIndex: %v\n", sol.ResultIdx)

}

func main() {

	println("Test 1: random Solution Generation")
	println("==================================")
	println("Example of random solution:")
	solution := verifiers.GenerateRandomSolution()
	printSolution(solution)

	println("\n\nTest 2: Solution from Index")
	println("==================================")
	success := true
	failueIndexes := []int{
		-1,
		-15,
		125,
		500,
		2000,
	}
	for _, v := range failueIndexes {
		s, err := verifiers.CalcSolutionFromIndex(v)
		if err == nil {
			fmt.Printf("Test Failed! %v should have failed, but it was a success. Solution: ", v)
			printSolution(s)
			success = false
		}
	}
	successIndexes := []int{}
	for i := range 125 {
		successIndexes = append(successIndexes, i)
	}
	for _, v := range successIndexes {
		s, err := verifiers.CalcSolutionFromIndex(v)
		if err != nil {
			fmt.Printf("Test Failed! %v should have succeeded, but it failed. Error:", v)
			fmt.Println(err)
			printSolution(s)
			success = false
		}
	}
	if success {
		println("Successfully passed the test!")
	}
}
