package main

import (
	"distasteful-bear/turing_machine/api"
	"distasteful-bear/turing_machine/run"
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
)

func main() {

	fmt.Println("Setting up Game...")

	// sol := verifiers.GenerateRandomSolution()

	sol := verifiers.Solution{ResultIdx: 0, Display: [3]rune{'1', '1', '1'}}

	puzzle := verifiers.GenerateRandomPuzzle(sol)

	// puzzle := verifiers.TestingVerifiers(verifiers.Solution{Idx: 0, Digits: [3]rune{'1', '1', '1'}})

	run.RunLocalSinglePlayer(puzzle)

	router := api.SetupRouter()
	// Start the server
	router.Run(":8080")

}
