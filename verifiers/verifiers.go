package verifiers

import (
	"fmt"
	"strconv"
	"strings"
)

type Solution struct {
	Idx    int // from 1-125
	Digits [3]rune
}

type Puzzle struct {
	Sol  Solution
	Vers [4]Verifier
}

type ResponseSet [125]bool
type Verifier struct {
	VerifierFunc func(sol Solution) bool
	Desc         string
}

// type VerifierBranch func() ResponseSet

type color string // blue, yellow, purple

func V_NumberOfNInSol(num rune, sol Solution) Verifier {
	numOfN := strings.Count(string(sol.Digits[:]), string(num))
	return Verifier{
		Desc: fmt.Sprintf("The number of %cs in the code", num),
		VerifierFunc: func(sol Solution) bool {
			return strings.Count(string([]rune(sol.Digits[:])), string(num)) == numOfN
		}}
}

func V_ColorCompareToNumber(col color, num rune, sol Solution) Verifier {
	var idxOfColor = 0
	switch col {
	case "blue":
		idxOfColor = 0
		break
	case "yellow":
		idxOfColor = 1
		break
	case "purple":
		idxOfColor = 2
		break
	default:
		panic("invalid color supplied to VerColorCompareToNumber")
	}

	v := Verifier{
		Desc: fmt.Sprintf("the %v number compared to %c", col, num),
	}

	if sol.Digits[idxOfColor] == num {
		v.VerifierFunc = func(sol Solution) bool {
			return sol.Digits[idxOfColor] == num
		}
	} else if sol.Digits[idxOfColor] > num {
		v.VerifierFunc = func(sol Solution) bool {
			return sol.Digits[idxOfColor] > num
		}
	} else if sol.Digits[idxOfColor] < num {
		v.VerifierFunc = func(sol Solution) bool {
			return sol.Digits[idxOfColor] < num
		}
	} else {
		panic("no equality evaluated to true in V_ColorCompareToNumber")
	}
	return v
}

func isTotalOdd(sol Solution) bool {
	total := 0
	for _, c := range sol.Digits {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			panic(fmt.Sprintf("Could not parse Int from Solution: %v", sol))
		}
		total += v
	}
	return (total % 2) != 0
}

func V_SumOfAllNumbersOddEven(sol Solution) Verifier {
	isOdd := isTotalOdd(sol)
	return Verifier{
		Desc: "if the sum of all the numbers is even or odd",
		VerifierFunc: func(sol Solution) bool {
			return isOdd == isTotalOdd(sol)
		},
	}
}

func countRepitions(sol Solution) int {
	maxRepititions := 0
	curRepitions := 0
	for i, c := range sol.Digits {
		nextChar := i + 1
		nextnextChar := i + 2
		if nextChar > 2 {
			nextChar -= 3
		}
		if nextnextChar > 2 {
			nextnextChar -= 3
		}
		if sol.Digits[nextChar] == c {
			curRepitions += 1
		}
		if sol.Digits[nextnextChar] == c {
			curRepitions += 1
		}
		if curRepitions > maxRepititions {
			maxRepititions = curRepitions
		}
		curRepitions = 0
	}
	return maxRepititions
}
func V_NumberRepeatsItself(sol Solution) Verifier {
	numRepetitions := countRepitions(sol)

	return Verifier{
		Desc: "how many times a number repeats itself in the code",
		VerifierFunc: func(sol Solution) bool {
			return numRepetitions == countRepitions(sol)
		},
	}

}

func TestingVerifiers(sol Solution) Puzzle {

	testVers := [4]Verifier{
		V_ColorCompareToNumber("yellow", '4', sol),
		V_NumberOfNInSol('3', sol),
		V_SumOfAllNumbersOddEven(sol),
		V_NumberRepeatsItself(sol),
	}

	return Puzzle{
		Sol:  sol,
		Vers: testVers,
	}
}
