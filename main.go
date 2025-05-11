package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"farkle/farkle"
)

const (
    ColorCyan   = "\033[36m"
    ColorReset  = "\033[0m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorBlue   = "\033[34m"
    ColorRed    = "\033[31m"
)

// CLI usage banner
const banner = ColorCyan + `=== Welcome to Farkle ===` + ColorReset + `
` + ColorGreen + `Single‑player:` + ColorReset + `
  ` + ColorYellow + `play [score]` + ColorReset + `                       → solo vs CPU
` + ColorGreen + `Multiplayer (2 players):` + ColorReset + `
  ` + ColorYellow + `play [score] --mp --create` + ColorReset + `         → host a lobby
  ` + ColorYellow + `play --mp --join=<ID> [--host=<ip>]` + ColorReset + `→ join a lobby (host defaults 127.0.0.1)
` + ColorBlue + `Score range 1000‑20000 (default 1000).` + ColorReset + `
` + ColorRed + `Type 'exit/quit' to quit.` + ColorReset

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(banner)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		tokens := strings.Fields(input)

		switch strings.ToLower(tokens[0]) {
		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		case "play":
			handlePlay(tokens[1:])

		default:
			fmt.Println("Unknown command. Use 'play' or 'exit'.")
		}
	}
}

// handlePlay parses flags & dispatches the correct game mode.
func handlePlay(args []string) {
	// Default parameters
	target := 1000
	isMP := false
	create := false
	joinID := ""
	hostIP := "127.0.0.1"

	for _, tok := range args {
		switch {
		case strings.HasPrefix(tok, "--"):
			switch {
			case tok == "--mp":
				isMP = true
			case tok == "--create":
				create = true
			case strings.HasPrefix(tok, "--join="):
				isMP = true
				joinID = strings.ToUpper(strings.TrimPrefix(tok, "--join="))
			case strings.HasPrefix(tok, "--host="):
				hostIP = strings.TrimPrefix(tok, "--host=")
			default:
				fmt.Println("Unknown flag:", tok)
				return
			}

		default: // treat as score
			if v, err := strconv.Atoi(tok); err == nil {
				if v < 1000 {
					v = 1000
				} else if v > 20000 {
					v = 20000
				}
				target = v
			} else {
				fmt.Println("Invalid score:", tok)
				return
			}
		}
	}

	// --- Dispatch ---
	if !isMP {
		farkle.WinningScore = target
		farkle.PlayGame()
		return
	}

	// Multiplayer validation
	if create && joinID != "" {
		fmt.Println("Cannot combine --create with --join.")
		return
	}
	if create {
		farkle.HostLobby(target)
		return
	}
	if joinID != "" {
		farkle.JoinLobby(hostIP, joinID)
		return
	}

	fmt.Println("For multiplayer, use --create OR --join=<ID>  [--host=<ip>].")
}
