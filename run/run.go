package run

import (
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"slices"
)

func printVerifiers(vs [4]verifiers.Verifier) {

	for i, v := range vs {
		fmt.Printf("Verifier #%v\n", (i + 1))
		fmt.Println(v)
	}
}
func enterGuess() [3]rune {
	for {
		fmt.Println("Please enter a three digit number with digits between 1-5:")
		var guess [3]rune
		_, err := fmt.Scanln(guess)
		invalid := false
		if err != nil {
			fmt.Println("failed to read input. please try again.")
		}
		if len(guess) != 3 {
			invalid = true
		}
		for _, c := range guess {
			if !slices.Contains([]rune{'1', '2', '3', '4', '5'}, c) {
				invalid = true
			}
		}
		if invalid {
			fmt.Println("input was invalid. please try again.")
		} else {
			return guess
		}
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
		fmt.Println("Guess: %v", guess)
	}
	return
}
