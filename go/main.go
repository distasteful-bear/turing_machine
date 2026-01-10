package main

import (
	"distasteful-bear/turing_machine/api"
	"distasteful-bear/turing_machine/run"
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"os"
	"slices"
)

func main() {

	fmt.Println("Setting up Game...")
	args := os.Args[1:]

	var sol verifiers.Solution
	if slices.Contains(args, "--known") {
		sol = verifiers.Solution{ResultIdx: 0, Display: [3]rune{'1', '1', '1'}}
	} else {
		sol = verifiers.GenerateRandomSolution()
	}
	if slices.Contains(args, "--local") {
		// cli game
		puzzle := verifiers.GenerateRandomPuzzle(sol)
		run.RunLocalSinglePlayer(puzzle)
	} else {
		// web api
		router := api.SetupRouter()
		router.Run(":8080")
	}
}
