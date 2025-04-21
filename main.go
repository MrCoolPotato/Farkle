package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"farkle/farkle"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("=== Welcome to Farkle ===")
	fmt.Println("Type 'play' to start a new game, or 'exit' to quit.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		switch strings.ToLower(input) {
		case "play":
			farkle.PlayGame()
		case "exit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type 'play' or 'exit'.")
		}
	}
}
