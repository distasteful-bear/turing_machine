package main

import (
	"distasteful-bear/turing_machine/run"
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"math/rand"
	"strconv"
)

func generateRandomSolution() solution {
	idx := rand.Intn(125)
	digits := strconv.FormatInt(int64(idx), 5)
	return solution{
		Idx:    idx,
		Digits: []rune(digits),
	}
}
func main() {

	fmt.Println("Setting up Game...")

	// numTests := 4
	// sol := generateRandomSolution()

	puzzle := verifiers.TestingVerifiers(verifiers.Solution{Idx: 0, Digits: []rune{'1', '1', '1'}})

	run.RunLocalSinglePlayer(puzzle)
}
