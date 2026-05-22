package run

import (
	"distasteful-bear/turing_machine/utils"
	"distasteful-bear/turing_machine/verifiers"
	"fmt"
	"strconv"
	"strings"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiBlue   = "\033[34m"
	ansiYellow = "\033[33m"
	ansiPurple = "\033[35m"
	ansiGreen  = "\033[32m"
	ansiRed    = "\033[31m"
)

func colorize(text, color string) string {
	return color + text + ansiReset
}

func colorName(name string) string {
	switch name {
	case "Blue", "blue":
		return colorize(name, ansiBlue)
	case "Yellow", "yellow":
		return colorize(name, ansiYellow)
	case "Purple", "purple":
		return colorize(name, ansiPurple)
	default:
		return name
	}
}

func colorizeVerifierDesc(desc string) string {
	desc = strings.ReplaceAll(desc, "blue", colorName("blue"))
	desc = strings.ReplaceAll(desc, "yellow", colorName("yellow"))
	desc = strings.ReplaceAll(desc, "purple", colorName("purple"))
	return desc
}

func formatSolution(sol verifiers.Solution) string {
	return fmt.Sprintf("%s%s%s",
		colorize(string(sol.Display[0]), ansiBlue),
		colorize(string(sol.Display[1]), ansiYellow),
		colorize(string(sol.Display[2]), ansiPurple),
	)
}

func printColorLegend() {
	fmt.Println("Digit colors / positions:")
	fmt.Printf("  %s digit = 1st digit in your guess\n", colorName("Blue"))
	fmt.Printf("  %s digit = 2nd digit in your guess\n", colorName("Yellow"))
	fmt.Printf("  %s digit = 3rd digit in your guess\n", colorName("Purple"))
	fmt.Println("Example: 421 means Blue=4, Yellow=2, Purple=1.")
}

func printVerifiers(displayOnly bool, vs [4]verifiers.VerifierBranch, guess verifiers.Solution) {
	for i, v := range vs {
		fmt.Printf("%sVerifier %v:%s ", ansiBold, (i + 1), ansiReset)
		fmt.Println(colorizeVerifierDesc(v.Desc))
		if !displayOnly {
			result := v.VerifierFunc(guess)
			fmt.Print("Result: ")
			if result {
				fmt.Println(colorize("true", ansiGreen))
			} else {
				fmt.Println(colorize("false", ansiRed))
			}
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
	fmt.Printf("%sWelcome to the Turing Machine Game.%s\n\n", ansiBold, ansiReset)

	fmt.Println("A Solution has been generated. Below are your verifiers.")
	fmt.Println("When you enter a guess, each verifier will return whether the test was successful.")
	fmt.Println("Use this information to deduce the solution.")
	fmt.Print("When you feel confident, you can end the guessing and check your work.\n\n")
	printColorLegend()
	fmt.Print("\n")

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
		printColorLegend()
		fmt.Print("\n")
		fmt.Printf("Guess: %s\n", formatSolution(sol))
		fmt.Print("\n")
		fmt.Printf("%s: %c\n", colorName("Blue"), sol.Display[0])
		fmt.Printf("%s: %c\n", colorName("Yellow"), sol.Display[1])
		fmt.Printf("%s: %c\n", colorName("Purple"), sol.Display[2])
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
			fmt.Printf("Guess: %s\n", formatSolution(p.Sol))
			fmt.Print("\n")
			fmt.Printf("%s: %c\n", colorName("Blue"), p.Sol.Display[0])
			fmt.Printf("%s: %c\n", colorName("Yellow"), p.Sol.Display[1])
			fmt.Printf("%s: %c\n", colorName("Purple"), p.Sol.Display[2])
			fmt.Print("\n\n")
			return
		}
	}
}
