package run

import (
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"slices"
	"strconv"
)

func printVerifiers(vs [4]verifiers.Verifier, guess verifiers.Solution) {
	for i, v := range vs {
		fmt.Printf("Verifier %v: ", (i + 1))
		fmt.Println(v.Desc)
		if [3]rune(guess.Digits) != [3]rune{'0', '0', '0'} {
			result := v.VerifierFunc(guess)
			fmt.Print("Result: ")
			fmt.Println(result)
		}
	}
}
func enterGuess() [3]rune {
	var finalGuess [3]rune
	for {
		fmt.Println("Please enter a three digit number with digits between 1-5:")
		var guess int
		_, err := fmt.Scan(&guess)
		invalid := false
		if err != nil {
			fmt.Println("failed to read input. ending process")
			break
		}
		strguess := strconv.Itoa(guess)
		if len(strguess) != 3 {
			continue
		}
		if len(strguess) != 3 {
			invalid = true
		}
		for _, c := range strguess {
			if !slices.Contains([]rune{'1', '2', '3', '4', '5'}, c) {
				invalid = true
			}
		}
		if invalid {
			fmt.Println("input was invalid. please try again.")
		} else {
			copy(finalGuess[:], []rune(strguess)[:])
			break
		}
	}
	return finalGuess

}
func RunLocalSinglePlayer(p verifiers.Puzzle) {
	fmt.Println("Welcome to the Turing Machine Game.\n\n")

	fmt.Println("A Solution has been generated. \nBelow are your verififers.")
	fmt.Println("When you enter a guess, each verifier will return whether the test was successful.")
	fmt.Println("Use this information to deduce the solution.")
	fmt.Println("When you feel confident, you can end the guessing and check your work.\n\n")

	checkAnswer := false

	intSolIdx, err := strconv.ParseInt("111", 5, 8)
	if err != nil {
		panic("could not parse int")
	}
	sol := verifiers.Solution{
		Digits: [3]rune{'6', '2', '1'},
		Idx:    int(intSolIdx),
	}
	for {
		if checkAnswer {
			break
		}
		printVerifiers(p.Vers, sol)

		fmt.Print("\n\n\n")
		guess := enterGuess()
		fmt.Printf("Guess: %c\n", guess)
		fmt.Printf("Blue: %c\n", guess[0])
		fmt.Printf("Yellow: %c\n", guess[1])
		fmt.Printf("Purple: %c\n", guess[2])

		guessIdx, err := strconv.ParseInt(string(guess[:]), 5, 8)
		if err != nil {
			panic("Failed to parse guess")
		}
		sol = verifiers.Solution{
			Digits: guess,
			Idx:    int(guessIdx),
		}
	}
}
