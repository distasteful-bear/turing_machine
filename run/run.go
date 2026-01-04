package run

import (
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
)

func printVerifiers(vs [4]verifiers.Verifier) {

	for i, v := range vs {
		fmt.Printf("Verifier #%v\n", (i + 1))
		fmt.Println(v)
	}
}
func RunLocalSinglePlayer(p verifiers.Puzzle) {
	fmt.Println("Welcome to the Turing Machine Game.\n\n")

	fmt.Println("A Solution has been generated. \nBelow are your verififers.")
	fmt.Println("When you enter a guess, each verifier will return whether the test was successful.")
	fmt.Println("Use this information to deduce the solution.")
	fmt.Println("When you feel confident, you can end the guessing and check your work.\n\n")

	checkAnswer := false
	for {
		if checkAnswer {
			break
		}
		printVerifiers(p.Vers)
		guess := enterGuess()

	}
	return
}
