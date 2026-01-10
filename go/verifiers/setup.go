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
	base5 := strconv.FormatInt(int64(idx), 5)

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
func getAllBranchesInVerifierForSolution(v Verifier, sol Solution) []VerifierBranch {
	allColors := GetColors()
	allDigits := GetDigitOptions()

	var results []VerifierBranch

	for _, c := range allColors {
		for _, d := range allDigits {
			results = append(results, v(c, d, sol))
		}
	}
	return results
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
func GenerateRandomPuzzle(sol Solution) Puzzle {

	vList := GetAllVerifiers()

	var allBranches []VerifierBranch
	for _, v := range vList {
		newBranches := getAllBranchesInVerifierForSolution(v, sol)
		allBranches = slices.Concat(allBranches, newBranches)
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
	innerLoopCount := 0
	outerLoopCount := 0
	for {
		outerLoopCount += 1
		for i, vr := range versAndResults {
			if slices.Contains(indexesOfPickedBranches, i) {
				continue
			}
			innerLoopCount += 1

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
		}
		if len(pickedBranches) == 4 {
			break
		}
		if outerLoopCount%1000 == 0 {
			indexesOfPickedBranches = []int{}
			pickedBranches = []VerifierBranch{}
			for i := range curSolCard {
				curSolCard[i] = true
			}
		}
		if outerLoopCount%1000 == 0 {
			fmt.Printf("Outer loop count: %v\n", outerLoopCount)
			fmt.Printf("Inner loop count: %v\n", innerLoopCount)
		}
	}
	return Puzzle{
		Sol:  sol,
		Vers: [4]VerifierBranch(pickedBranches),
	}
}
func AssembleTestingVerifiers(sol Solution) Puzzle {

	testVers := [4]VerifierBranch{
		V_ColorCompareToNumber("yellow", '4', sol),
		V_NumberOfNInSol("yellow", '3', sol),
		V_SumOfAllNumbersOddEven("yellow", '3', sol),
		V_NumberRepeatsItself("yellow", '3', sol),
	}

	return Puzzle{
		Sol:  sol,
		Vers: testVers,
	}
}
