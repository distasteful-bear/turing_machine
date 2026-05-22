package verifiers

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
)

type resultSet [125]bool

type BranchAndResult struct {
	Branch VerifierBranch
	Result resultSet
}

func GenerateRandomSolution() Solution {
	idx := rand.Intn(125)
	solution, err := CalcSolutionFromIndex(idx)
	if err != nil {
		// errors should never happen, but you can just reroll.
		return GenerateRandomSolution()
	}
	return solution
}
func CalcSolutionFromIndex(idx int) (Solution, error) {
	if idx < 0 || idx >= 125 {
		return Solution{}, error(fmt.Errorf("invalid index for solution!!! %v", idx))
	}

	base5 := fmt.Sprintf("%03s", strconv.FormatInt(int64(idx), 5))

	var display = [3]rune{'1', '1', '1'}
	for i, c := range base5 {
		display[i] = (c + 1)
	}

	return Solution{
		ResultIdx: idx,
		Display:   display,
	}, nil
}
func GetColors() []color {
	return []color{
		"blue",
		"yellow",
		"purple",
	}
}
func GetDigitOptions() [5]rune {
	return [5]rune{
		'1',
		'2',
		'3',
		'4',
		'5',
	}
}
func countTrue(rs resultSet) int {
	ttl := 0
	for _, b := range rs {
		if b {
			ttl += 1
		}
	}
	return ttl
}
func generateResultCardForBranch(vb VerifierBranch) resultSet {
	var resultSet [125]bool
	for i := range 125 {
		s, err := CalcSolutionFromIndex(i)
		if err != nil {
			panic("invalid index generated for result card generation!")
		}
		resultSet[i] = vb.VerifierFunc(s)
	}
	return resultSet
}
func intersectionOfSets(rsList []resultSet) resultSet {
	if len(rsList) == 0 {
		panic("no sets provided!")
	}
	if len(rsList) == 1 {
		return rsList[0]
	}

	var final resultSet = rsList[0]
	for i, rs := range rsList {
		if i == 0 {
			continue // final is already set to the first result set
		}
		for k, b := range rs {
			if b == false {
				final[k] = false
			}
		}
	}
	return final
}
func getRandomVerifier(vrList []BranchAndResult) (int, BranchAndResult) {
	randomIndex := rand.Intn(len(vrList))
	return randomIndex, vrList[randomIndex]
}

func GenerateRandomPuzzle() Puzzle {

	var sol = GenerateRandomSolution()
	var allBranches []VerifierBranch

	allIterators := GetAllVerifiers()

	for _, v := range allIterators.Solution {
		allBranches = append(allBranches, v(sol))
	}
	for _, v := range allIterators.Color {
		for _, c := range GetColors() {
			allBranches = append(allBranches, v(c, sol))
		}
	}
	for _, v := range allIterators.Number {
		for _, n := range GetDigitOptions() {
			allBranches = append(allBranches, v(n, sol))
		}
	}
	for _, v := range allIterators.ColorColor {
		for _, c1 := range GetColors() {
			for _, c2 := range GetColors() {
				allBranches = append(allBranches, v(c1, c2, sol))
			}
		}
	}
	for _, v := range allIterators.ColorNumber {
		for _, c := range GetColors() {
			for _, n := range GetDigitOptions() {
				allBranches = append(allBranches, v(c, n, sol))
			}
		}
	}

	var versAndResults []BranchAndResult
	for _, b := range allBranches {
		result := generateResultCardForBranch(b)
		versAndResults = append(versAndResults, BranchAndResult{
			Branch: b,
			Result: result,
		})
	}

	var indexesOfPickedBranches []int
	var pickedBranches []VerifierBranch
	var curSolCard [125]bool
	for i := range curSolCard {
		curSolCard[i] = true
	}
	outerLoopCount := 0
	for {
		outerLoopCount += 1
		i, vr := getRandomVerifier(versAndResults)
		if slices.Contains(indexesOfPickedBranches, i) {
			continue
		}

		solutionCardIfPicked := intersectionOfSets([]resultSet{vr.Result, curSolCard})
		before := countTrue(curSolCard)
		after := countTrue(solutionCardIfPicked)

		if len(pickedBranches) < 3 {
			if before > after && after != 1 {
				indexesOfPickedBranches = append(indexesOfPickedBranches, i)
				pickedBranches = append(pickedBranches, vr.Branch)
				curSolCard = solutionCardIfPicked
			}
		} else {
			if after == 1 {
				indexesOfPickedBranches = append(indexesOfPickedBranches, i)
				pickedBranches = append(pickedBranches, vr.Branch)
				curSolCard = solutionCardIfPicked
				break
			}
		}
		if len(pickedBranches) == 4 {
			break
		}
		if outerLoopCount%10000 == 0 {
			return GenerateRandomPuzzle()
		}
		if outerLoopCount%10000 == 0 {
			fmt.Printf("Outer loop count: %v\n", outerLoopCount)
		}
	}
	return Puzzle{
		Sol:  sol,
		Vers: [4]VerifierBranch(pickedBranches),
	}
}
