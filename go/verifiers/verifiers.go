package verifiers

import (
	"fmt"
	"strconv"
	"strings"
)

type Solution struct {
	ResultIdx int     // from 0-124
	Display   [3]rune // 1-5
}

type Puzzle struct {
	Sol  Solution
	Vers [4]VerifierBranch
}

type ResponseSet [125]bool

type VerifierBranch struct {
	VerifierFunc func(sol Solution) bool
	Desc         string
}
type Verifier func(color, rune, Solution) VerifierBranch

type color string // blue, yellow, purple
type digit rune   // 1,2,3,4, or 5

func V_NumberOfNInSol(col color, num rune, sol Solution) VerifierBranch {
	// intentionally do nothing with color
	numOfN := strings.Count(string(sol.Display[:]), string(num))
	return VerifierBranch{
		Desc: fmt.Sprintf("The number of %cs in the code", num),
		VerifierFunc: func(sol Solution) bool {
			return strings.Count(string([]rune(sol.Display[:])), string(num)) == numOfN
		}}
}

func V_ColorCompareToNumber(col color, num rune, sol Solution) VerifierBranch {
	var idxOfColor = 0
	switch col {
	case "blue":
		idxOfColor = 0
	case "yellow":
		idxOfColor = 1
	case "purple":
		idxOfColor = 2
	default:
		panic("invalid color supplied to VerColorCompareToNumber")
	}

	v := VerifierBranch{
		Desc: fmt.Sprintf("the %v number compared to %c", col, num),
	}

	if sol.Display[idxOfColor] == num {
		v.VerifierFunc = func(sol Solution) bool {
			return sol.Display[idxOfColor] == num
		}
	} else if sol.Display[idxOfColor] > num {
		v.VerifierFunc = func(sol Solution) bool {
			return sol.Display[idxOfColor] > num
		}
	} else if sol.Display[idxOfColor] < num {
		v.VerifierFunc = func(sol Solution) bool {
			return sol.Display[idxOfColor] < num
		}
	} else {
		panic("no equality evaluated to true in V_ColorCompareToNumber")
	}
	return v
}

func isTotalOdd(sol Solution) bool {
	total := 0
	for _, c := range sol.Display {
		v, err := strconv.Atoi(string(c))
		if err != nil {
			panic(fmt.Sprintf("Could not parse Int from Solution: %v", sol))
		}
		total += v
	}
	return (total % 2) != 0
}
func V_SumOfAllNumbersOddEven(col color, num rune, sol Solution) VerifierBranch {
	isOdd := isTotalOdd(sol)
	return VerifierBranch{
		Desc: "if the sum of all the numbers is even or odd",
		VerifierFunc: func(sol Solution) bool {
			return isOdd == isTotalOdd(sol)
		},
	}
}

func countRepitions(sol Solution) int {
	maxRepititions := 0
	curRepitions := 0
	for i, c := range sol.Display {
		nextChar := i + 1
		nextnextChar := i + 2
		if nextChar > 2 {
			nextChar -= 3
		}
		if nextnextChar > 2 {
			nextnextChar -= 3
		}
		if sol.Display[nextChar] == c {
			curRepitions += 1
		}
		if sol.Display[nextnextChar] == c {
			curRepitions += 1
		}
		if curRepitions > maxRepititions {
			maxRepititions = curRepitions
		}
		curRepitions = 0
	}
	return maxRepititions
}
func V_NumberRepeatsItself(col color, num rune, sol Solution) VerifierBranch {
	numRepetitions := countRepitions(sol)

	return VerifierBranch{
		Desc: "how many times a number repeats itself in the code",
		VerifierFunc: func(sol Solution) bool {
			return numRepetitions == countRepitions(sol)
		},
	}
}

func GetAllVerifiers() []Verifier {
	AllVerifiers := []Verifier{
		V_NumberOfNInSol,
		V_ColorCompareToNumber,
		V_SumOfAllNumbersOddEven,
		V_NumberRepeatsItself,
	}
	return AllVerifiers
}
