package verifiers

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func validDisplaySolutions(p Puzzle) []string {
	var valid []string

	for blue := '1'; blue <= '5'; blue++ {
		for yellow := '1'; yellow <= '5'; yellow++ {
			for purple := '1'; purple <= '5'; purple++ {
				sol := Solution{
					Display: [3]rune{blue, yellow, purple},
				}

				passesAll := true
				for _, verifier := range p.Vers {
					if !verifier.VerifierFunc(sol) {
						passesAll = false
						break
					}
				}

				if passesAll {
					valid = append(valid, string(sol.Display[:]))
				}
			}
		}
	}

	return valid
}

func TestGeneratedPuzzlesHaveSingleValidDisplaySolution(t *testing.T) {
	samples := 10_000
	if override := os.Getenv("PUZZLE_UNIQUENESS_SAMPLES"); override != "" {
		parsed, err := strconv.Atoi(override)
		if err != nil || parsed < 1 {
			t.Fatalf("invalid PUZZLE_UNIQUENESS_SAMPLES value %q", override)
		}
		samples = parsed
	}

	rand.Seed(1)

	for i := 0; i < samples; i++ {
		puzzle := GenerateRandomPuzzle()
		validSolutions := validDisplaySolutions(puzzle)

		if len(validSolutions) != 1 {
			descriptions := make([]string, 0, len(puzzle.Vers))
			for _, verifier := range puzzle.Vers {
				descriptions = append(descriptions, verifier.Desc)
			}

			t.Fatalf(
				"generated puzzle %d has %d valid display solutions; stored solution %q; valid solutions %v; verifiers %v",
				i+1,
				len(validSolutions),
				string(puzzle.Sol.Display[:]),
				validSolutions,
				descriptions,
			)
		}

		if validSolutions[0] != string(puzzle.Sol.Display[:]) {
			t.Fatalf(
				"generated puzzle %d unique valid solution %q does not match stored solution %q",
				i+1,
				validSolutions[0],
				string(puzzle.Sol.Display[:]),
			)
		}
	}
}
