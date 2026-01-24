package run

import (
	"distasteful-bear/turing_machine/utils"
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"strconv"
)

func printVerifiers(displayOnly bool, vs [4]verifiers.VerifierBranch, guess verifiers.Solution) {
	for i, v := range vs {
		fmt.Printf("Verifier %v: ", (i + 1))
		fmt.Println(v.Desc)
		if !displayOnly {
			result := v.VerifierFunc(guess)
			fmt.Print("Result: ")
			fmt.Println(result)
		}
	}
}

func enterGuess() verifiers.Solution {
	// returns within loop, continue whenever an issue occurs
	for {
		fmt.Println("Please enter a three digit number with digits between 1-5:")
		// get int from scan
		var guess int
		_, err := fmt.Scan(&guess)
		if err != nil {
			fmt.Println("failed to read input, scaning to int failed. try again.")
			continue
		}
		// only three chars
		strguess := strconv.Itoa(guess)
		sol, err := utils.SanitizeGuess(strguess)
		if err != nil {
			fmt.Println(err)
			continue
		}
		return sol
	}
}

func RunLocalSinglePlayer(p verifiers.Puzzle) {
	utils.CallClear()
	fmt.Println("Welcome to the Turing Machine Game.\n\n")

	fmt.Println("A Solution has been generated. \nBelow are your verififers.")
	fmt.Println("When you enter a guess, each verifier will return whether the test was successful.")
	fmt.Println("Use this information to deduce the solution.")
	fmt.Println("When you feel confident, you can end the guessing and check your work.\n\n")

	checkAnswer := false

	intSolIdx, err := strconv.ParseInt("421", 5, 8)
	if err != nil {
		panic("could not parse int")
	}
	sol := verifiers.Solution{
		Display:   [3]rune{'4', '2', '1'},
		ResultIdx: int(intSolIdx),
	}
	for {
		if checkAnswer {
			break
		}

		// enter guesses
		printVerifiers(true, p.Vers, sol)
		fmt.Print("\n\n\n")
		sol := enterGuess()

		// display results
		utils.CallClear()
		fmt.Printf("Guess: %c\n", sol.Display)
		fmt.Print("\n")
		fmt.Printf("Blue: %c\n", sol.Display[0])
		fmt.Printf("Yellow: %c\n", sol.Display[1])
		fmt.Printf("Purple: %c\n", sol.Display[2])
		fmt.Print("\n\n")
		printVerifiers(false, p.Vers, sol)

		fmt.Print("\n\n")

		fmt.Println("Are you ready to guess the final solution? \n(0 = no, 1 = yes)")
		var readyToCheck int
		_, err = fmt.Scan(&readyToCheck)
		if err != nil {
			fmt.Println("failed to read input. assuming you would like to continue guessing,")
			continue
		}
		if readyToCheck == 1 {
			checkAnswer = true
		}
	}

	for {
		sol := enterGuess()

		if sol.ResultIdx == p.Sol.ResultIdx {
			fmt.Println("Success! You have solved the puzzle.")
			return
		} else {
			fmt.Println("Failure! You have not solved the puzzle.")
			fmt.Println("The answer was:")
			fmt.Printf("Guess: %c\n", sol.Display)
			fmt.Print("\n")
			fmt.Printf("Blue: %c\n", sol.Display[0])
			fmt.Printf("Yellow: %c\n", sol.Display[1])
			fmt.Printf("Purple: %c\n", sol.Display[2])
			fmt.Print("\n\n")
			return
		}
	}
}
