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

const winningScore = 10000

type Player struct {
	Name  string
	Total int
}

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

	player := Player{Name: "You"}
	enemy := Player{Name: "Enemy"}

	round := 1
	for {
		fmt.Printf("\n========================\n")
		fmt.Printf(" ROUND %d – First to %d\n", round, winningScore)
		fmt.Printf("========================\n")
		fmt.Printf("Scoreboard → You: %d | Enemy: %d\n", player.Total, enemy.Total)

		fmt.Println("\nYour turn:")
		fmt.Println("First roll will happen automatically; then choose dice to keep or bank.")
		playerPoints := playerTurn()
		player.Total += playerPoints
		fmt.Printf("You banked %d points. New total: %d\n", playerPoints, player.Total)
		if player.Total >= winningScore {
			fmt.Println("\n🏆  You win the game! 🏆")
			return
		}

		fmt.Println("\nEnemy turn:")
		enemyPoints := enemyTurn()
		enemy.Total += enemyPoints
		fmt.Printf("Enemy banked %d points. New total: %d\n", enemyPoints, enemy.Total)
		if enemy.Total >= winningScore {
			fmt.Println("\n💀  Enemy wins the game. Better luck next time!")
			return
		}

		round++
	}
}

func rollDice(n int) []int {
	dice := make([]int, n)
	for i := 0; i < n; i++ {
		dice[i] = rand.Intn(6) + 1
	}
	return dice
}

func renderDice(dice []int) {
	for _, d := range dice {
		fmt.Printf("[%s %d] ", dieFaces[d], d)
	}
	fmt.Println()
}

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

func promptAction(roll []int) (kept []int, action string) {
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

		parseDice := func(parts []string) ([]int, bool) {
			req := make(map[int]int)
			var vals []int
			for _, p := range parts {
				val, err := strconv.Atoi(p)
				if err != nil || val < 1 || val > 6 {
					fmt.Println("Invalid die value:", p)
					return nil, false
				}
				req[val]++
				vals = append(vals, val)
			}
			avail := make(map[int]int)
			for _, d := range roll {
				avail[d]++
			}
			for val, cnt := range req {
				if cnt > avail[val] {
					fmt.Printf("Cannot keep %d of '%d'; only %d available.\n", cnt, val, avail[val])
					return nil, false
				}
			}
			for _, v := range vals {
				if calculateScore([]int{v}) == 0 {
					fmt.Printf("Die %d does not score; only scoring dice may be kept.\n", v)
					return nil, false
				}
			}
			return vals, true
		}

		if strings.HasPrefix(lower, "keep ") {
			parts := strings.Fields(input)[1:]
			vals, ok := parseDice(parts)
			if ok {
				return vals, "keep"
			}
			continue
		}

		if strings.HasPrefix(lower, "bank") {
			parts := strings.Fields(input)[1:]
			if len(parts) == 0 {
				fmt.Println("Specify dice to bank, e.g., 'bank 1 5'.")
				continue
			}
			vals, ok := parseDice(parts)
			if ok {
				return vals, "bank"
			}
			continue
		}

		fmt.Println("Commands: 'keep X X...' (score & continue), 'bank X X...' (score & pass), or 'quit'")
	}
}

func playerTurn() int {
	diceToRoll := 6
	turnScore := 0

outer:
	for {
		fmt.Printf("-- Rolling %d dice --\n", diceToRoll)
		roll := rollDice(diceToRoll)
		renderDice(roll)

		if calculateScore(roll) == 0 {
			fmt.Println("Farkle! You lose all unbanked points for this turn.")
			return 0
		}

		fmt.Println("Commands: 'keep X X...' to score & CONTINUE, 'bank X X...' to score & PASS, or 'quit'")

	inner:
		for {
			kept, action := promptAction(roll)
			switch action {
			case "quit":
				fmt.Println("Goodbye!")
				os.Exit(0)
			case "keep":
				score := calculateScore(kept)
				turnScore += score
				fmt.Printf("Scored %d (turn total %d). Continuing...\n", score, turnScore)

				diceToRoll -= len(kept)
				if diceToRoll == 0 {
					fmt.Println("Hot dice! All dice scored, rolling 6 fresh dice.")
					diceToRoll = 6
				}
				break inner
			case "bank":
				score := calculateScore(kept)
				turnScore += score
				return turnScore
			}
		}
		continue outer
	}
}

func enemyTurn() int {
	roll := rollDice(6)
	renderDice(roll)
	score := calculateScore(roll)
	if score == 0 {
		fmt.Println("Enemy Farkled and scores 0.")
		return 0
	}
	fmt.Printf("Enemy keeps all scoring dice and banks %d.\n", score)
	return score
}
