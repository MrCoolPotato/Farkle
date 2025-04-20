package farkle

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var dieFaces = map[int]string{
	1: "⚀",
	2: "⚁",
	3: "⚂",
	4: "⚃",
	5: "⚄",
	6: "⚅",
}

func PlayGame() {
	rand.Seed(time.Now().UnixNano())
	round := 1

	for {
		fmt.Printf("\n== Round %d ==\n", round)
		// Player's turn
		diceToRoll := 6
		var keptValues []int
		var turnScore int

	playerLoop:
		for {
			fmt.Printf("-- Player Roll: rolling %d dice --\n", diceToRoll)
			dice := rollDice(diceToRoll)
			renderDice(dice)
			fmt.Printf("Kept so far: ")
			renderDice(keptValues)
			fmt.Printf("Potential score: %d\n", calculateScore(keptValues))
			fmt.Println("Commands: 'keep X X...', 'roll' to score & reroll, 'bank' to score & end turn")

			kept, action := promptAction(dice)
			switch action {
			case "quit":
				fmt.Println("Goodbye!")
				os.Exit(0)
			case "keep":
				keptValues = append(keptValues, kept...)
				diceToRoll -= len(kept)
				if diceToRoll <= 0 {
					fmt.Println("Hot dice! Resetting to 6 dice.")
					diceToRoll = 6
					keptValues = []int{}
				}
				continue
			case "roll":
				if len(keptValues) == 0 {
					fmt.Println("You must keep at least one die before rolling.")
					continue
				}
				// score & reroll: calculate but continue
				turnScore = calculateScore(keptValues)
				break playerLoop
			case "bank":
				turnScore = calculateScore(keptValues)
				break playerLoop
			}
		}

		fmt.Printf("You scored %d points this turn.\n", turnScore)

		// Enemy's turn: single roll (no scoring for now)
		fmt.Println("-- Enemy's turn --")
		enemyDice := rollDice(6)
		renderDice(enemyDice)

		round++
		break // adjust to allow multiple rounds later
	}
}

// rollDice generates n random dice values (1–6).
func rollDice(n int) []int {
	dice := make([]int, n)
	for i := 0; i < n; i++ {
		dice[i] = rand.Intn(6) + 1
	}
	return dice
}

// renderDice displays each die with its Unicode face and numeric value.
func renderDice(dice []int) {
	for _, d := range dice {
		fmt.Printf("[%s %d] ", dieFaces[d], d)
	}
	fmt.Println()
}

// calculateScore computes Farkle points for the given kept dice.
func calculateScore(dice []int) int {
	counts := make(map[int]int)
	for _, d := range dice {
		counts[d]++
	}
	score := 0
	for val, cnt := range counts {
		if cnt >= 3 {
			if val == 1 {
				score += 1000
			} else {
				score += val * 100
			}
			cnt -= 3
		}
		if val == 1 {
			score += cnt * 100
		} else if val == 5 {
			score += cnt * 50
		}
	}
	return score
}

// promptAction parses user input for keep, roll, bank, or quit commands.
func promptAction(dice []int) (kept []int, action string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			os.Exit(0)
		}
		input := strings.TrimSpace(scanner.Text())
		lower := strings.ToLower(input)
		if lower == "quit" || lower == "exit" {
			return nil, "quit"
		}
		if lower == "bank" {
			return nil, "bank"
		}
		if lower == "roll" {
			return nil, "roll"
		}
		if strings.HasPrefix(lower, "keep ") {
			parts := strings.Fields(strings.TrimPrefix(input, "keep "))
			req := make(map[int]int)
			valid := true
			for _, p := range parts {
				val, err := strconv.Atoi(p)
				if err != nil || val < 1 || val > 6 {
					fmt.Println("Invalid die value:", p)
					valid = false
					break
				}
				req[val]++
			}
			if !valid {
				continue
			}
			// Check availability
			avail := make(map[int]int)
			for _, d := range dice {
				avail[d]++
			}
			for val, cnt := range req {
				if cnt > avail[val] {
					fmt.Printf("Cannot keep %d of '%d'; only %d available.\n", cnt, val, avail[val])
					valid = false
					break
				}
			}
			if !valid {
				continue
			}
			// Build kept slice
			for _, p := range parts {
				val, _ := strconv.Atoi(p)
				kept = append(kept, val)
			}
			return kept, "keep"
		}
		fmt.Println("Invalid command. Use 'keep X...', 'roll', or 'bank'.")
	}
}
