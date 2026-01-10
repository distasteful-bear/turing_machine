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
	if slices.Contains(args, "--local") {
		// cli game
		puzzle := verifiers.GenerateRandomPuzzle()
		run.RunLocalSinglePlayer(puzzle)
	} else {
		// web api
		router := api.SetupRouter()
		router.Run(":8080")
	}
}
