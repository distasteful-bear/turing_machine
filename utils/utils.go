package utils

import (
	"distasteful-bear/turing_machine/verifiers"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strconv"
)

func CallClear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows: use "cmd /c cls"
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		// Unix-like systems (Linux, macOS, etc.): use "clear"
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func SanitizeGuess(guess string) (verifiers.Solution, error) {
	solution := verifiers.Solution{
		ResultIdx: 0,
		Display:   [3]rune{'1', '1', '1'},
	}
	if len(guess) != 3 {
		return solution, errors.New("guess must be exactly 3 characters long")
	}
	for _, c := range guess {
		if !slices.Contains([]rune{'1', '2', '3', '4', '5'}, c) {
			return solution, errors.New("guess contains invalid characters. Only use the numbers 1 through 5")
		}
	}
	runesForBase5 := []rune(guess)
	for i, c := range runesForBase5 {
		runesForBase5[i] = c - 1 // need to ensure this rune string works with base 5 conversion which is 0->4
	}

	guessIdx, err := strconv.ParseInt(string(runesForBase5[:]), 5, 8)
	if err != nil {
		return solution, errors.New("guess contains non-numeric characters")
	}

	runeSlice := []rune(guess)
	solution.Display = [3]rune{runeSlice[0], runeSlice[1], runeSlice[2]}
	solution.ResultIdx = int(guessIdx)
	fmt.Printf("Guess: %s, Index: %d, Display: %c%c%c\n", guess, guessIdx, solution.Display[0], solution.Display[1], solution.Display[2])
	return solution, nil
}
