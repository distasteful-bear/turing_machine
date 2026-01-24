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

type SolutionVerifier func(Solution) VerifierBranch
type ColorVerifier func(color, Solution) VerifierBranch
type NumberVerifier func(rune, Solution) VerifierBranch
type ColorColorVerifier func(color, color, Solution) VerifierBranch
type ColorNumberVerifier func(color, rune, Solution) VerifierBranch

type VerifierBranch struct {
	VerifierFunc func(sol Solution) bool
	Desc         string
}

type color string // blue, yellow, purple
type digit rune   // 1,2,3,4, or 5

// ColorNumberVerifiers
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

// NumberVerifiers
func V_NumberOfNInSol(num rune, sol Solution) VerifierBranch {
	// intentionally do nothing with color
	numOfN := strings.Count(string(sol.Display[:]), string(num))
	return VerifierBranch{
		Desc: fmt.Sprintf("The number of %cs in the code", num),
		VerifierFunc: func(sol Solution) bool {
			return strings.Count(string([]rune(sol.Display[:])), string(num)) == numOfN
		}}
}

// SolutionVerifiers
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
func V_SumOfAllNumbersOddEven(sol Solution) VerifierBranch {
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
func V_NumberRepeatsItself(sol Solution) VerifierBranch {
	numRepetitions := countRepitions(sol)

	return VerifierBranch{
		Desc: "how many times a number repeats itself in the code",
		VerifierFunc: func(sol Solution) bool {
			return numRepetitions == countRepitions(sol)
		},
	}
}

// ColorVerifiers
func getNumFromColor(c color, sol Solution) int {
	switch c {
	case "blue":
		return 0
	case "yellow":
		return 1
	case "purple":
		return 2
	default:
		panic("invalid color")
	}
}

func smallestColor(sol Solution) bool {
	blue := 0
	yellow := 1
	purple := 2
	if sol.Display[blue] < sol.Display[yellow] &&
			sol.Display[blue] < sol.Display[purple] {
				return "blue"
			}
	case "yellow":
		return sol.Display[1] < sol.Display[getNumFromColor("blue", sol)] &&
			sol.Display[1] < sol.Display[getNumFromColor("purple", sol)]
	case "purple":
		return sol.Display[2] < sol.Display[getNumFromColor("blue", sol)] &&
			sol.Display[2] < sol.Display[getNumFromColor("yellow", sol)]
	default:
		panic("invalid color")
	}
}
func V_ColorIsSmallest(c color, sol Solution) VerifierBranch {
	solAns := isColorSmallest(c, sol)
	switch c {
	case "blue":
		return VerifierBranch{
			Desc: "is blue less than or equal to all other numbers",
			VerifierFunc: func(sol Solution) bool {
				blueIdx := 0
				return solAns == isColorSmallest(sol)
			},
		}
	case "yellow":
		return VerifierBranch{
			Desc: "is yellow less than or equal to all other numbers",
			VerifierFunc: func(sol Solution) bool {
				yellowIdx := 1
				return sol.Display[yellowIdx] < sol.Display[getNumFromColor("blue", sol)] &&
					sol.Display[yellowIdx] < sol.Display[getNumFromColor("purple", sol)]
			},
		}
	case "purple":
		return VerifierBranch{
			Desc: "is purple less than or equal to all other numbers",
			VerifierFunc: func(sol Solution) bool {
				purpleIdx := 2
				return sol.Display[purpleIdx] < sol.Display[getNumFromColor("blue", sol)] &&
					sol.Display[purpleIdx] < sol.Display[getNumFromColor("yellow", sol)]
			},
		}
	default:
		panic("invalid color")
	}
}
func V_ColorGreaterThanColor(c1 color, c2 color, sol Solution) VerifierBranch {
	return VerifierBranch{
		Desc: fmt.Sprintf("%s > %s", c1, c2),
		VerifierFunc: func(sol Solution) bool {
			return sol.Display[getNumFromColor(c1, sol)] > sol.Display[getNumFromColor(c2, sol)]
		},
	}
}

type VerifierIterators struct {
	Solution    []SolutionVerifier
	Color       []ColorVerifier
	Number      []NumberVerifier
	ColorColor  []ColorColorVerifier
	ColorNumber []ColorNumberVerifier
}

func GetAllVerifiers() VerifierIterators {
	return VerifierIterators{
		Solution:    []SolutionVerifier{V_NumberRepeatsItself, V_SumOfAllNumbersOddEven},
		Color:       []ColorVerifier{V_ColorIsSmallest},
		Number:      []NumberVerifier{V_NumberOfNInSol},
		ColorColor:  []ColorColorVerifier{},
		ColorNumber: []ColorNumberVerifier{V_ColorCompareToNumber},
	}
}
