package main

import (
	"distasteful-bear/turing_machine/api"
	"distasteful-bear/turing_machine/run"
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file for Auth0 configuration
	if err := godotenv.Load("api/.env"); err != nil {
		log.Println("Warning: Could not load api/.env file:", err)
	}

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
