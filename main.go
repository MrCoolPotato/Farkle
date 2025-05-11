package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"farkle/farkle"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("=== Welcome to Farkle ===")
	fmt.Println("Type 'play [score]' to start a new game (min 1000, max 20000).")
	fmt.Println("Leave the score empty for 1000. Type 'exit' to quit.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(strings.ToLower(input))
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "play":
			target := 1000 // default
			// Sanity check for score
			if len(parts) > 1 {
				if v, err := strconv.Atoi(parts[1]); err == nil {
					if v < 1000 {
						v = 1000
					} else if v > 20000 {
						v = 20000
					}
					target = v
				}
			}
			farkle.WinningScore = target
			farkle.PlayGame()

		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Unknown command. Type 'play [score]' or 'exit'.")
		}
	}
}
