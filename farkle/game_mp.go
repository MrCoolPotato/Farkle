package farkle

import (
	crand "crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

//MARK: Message Definition
type NetMsg struct {
	T         string `json:"t"`
	Dice      []int  `json:"dice,omitempty"`
	Keep      []int  `json:"keep,omitempty"`
	Bank      bool   `json:"bank,omitempty"`
	Idx       int    `json:"idx,omitempty"`
	Delta     int    `json:"delta,omitempty"`
	Total     int    `json:"total,omitempty"`
	Target    int    `json:"target,omitempty"`
	Name      string `json:"name,omitempty"`
	Round     int    `json:"round,omitempty"`
	HostTotal int    `json:"htotal,omitempty"`
	PeerTotal int    `json:"ptotal,omitempty"`
}

//MARK: Host Lobby
func HostLobby(target int) {
	WinningScore = target
	id := generateLobbyID()
	fmt.Println(ColorYellow+"Lobby created. Share ID: "+id+ColorReset)

	ln, err := net.Listen("tcp", ":9313")
	if err != nil {
		fmt.Println(ColorRed+"Listen error:", err, ColorReset)
		return
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(ColorRed+"Accept error:", err, ColorReset)
		return
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	var hello NetMsg
	if err := dec.Decode(&hello); err != nil || hello.T != "hello" {
		fmt.Println(ColorRed+"Handshake failed", ColorReset)
		return
	}
	enc.Encode(NetMsg{T: "welcome", Idx: 1, Target: WinningScore})

	hostTotal, peerTotal := 0, 0
	currentIdx := 0
	diceToRoll := 6
	encHot := make(chan NetMsg, 2)
	go func() { for m := range encHot { enc.Encode(m) } }()
	round := 1

	for {
		if currentIdx == 0 && diceToRoll == 6 {
			fmt.Printf("\n========================\n")
			fmt.Printf(" ROUND %d – First to %d\n", round, WinningScore)
			fmt.Printf("========================\n")
			fmt.Printf("Scoreboard → %sYou%s: %d | %sPeer%s: %d\n", ColorGreen, ColorReset, hostTotal, ColorRed, ColorReset, peerTotal)
			enc.Encode(NetMsg{T: "banner", Round: round, HostTotal: hostTotal, PeerTotal: peerTotal, Target: WinningScore})
			round++
			fmt.Println("\n" + ColorGreen + "Your turn:" + ColorReset)
			fmt.Println(ColorYellow + "First roll will happen automatically; then choose dice to keep or bank." + ColorReset)
		}

		roll := rollDice(diceToRoll)
		enc.Encode(NetMsg{T: "roll", Dice: roll, Idx: currentIdx})
		if currentIdx == 0 {
			renderDice(roll)
		}

		if calculateScore(roll) == 0 {
			enc.Encode(NetMsg{T: "farkle", Idx: currentIdx})
			if currentIdx == 0 {
				fmt.Println(ColorRed + "Farkle! You scored 0 this turn." + ColorReset)
			}
			diceToRoll = 6
			currentIdx = 1 - currentIdx
			continue
		}

		if currentIdx == 0 {
			turnScore, diceRemaining, ended := hostTurnLoop(roll, diceToRoll, encHot)
			hostTotal += turnScore
			enc.Encode(NetMsg{T: "score", Idx: 0, Delta: turnScore, Total: hostTotal})
			if hostTotal >= WinningScore {
				enc.Encode(NetMsg{T: "game_over", Idx: 0})
				fmt.Println(ColorGreen + "🏆 You win! Returning to menu." + ColorReset)
				return
			}
			diceToRoll = diceRemaining
			if !ended {
				continue
			}
		} else {
			enc.Encode(NetMsg{T: "your_turn"})
			var act NetMsg
			if err := dec.Decode(&act); err != nil || act.T != "action" {
				fmt.Println(ColorRed + "Peer disconnected." + ColorReset)
				return
			}
			score := calculateScore(act.Keep)
			if score == 0 {
				enc.Encode(NetMsg{T: "farkle", Idx: 1})
				diceToRoll = 6
				fmt.Println(ColorRed + "Peer Farkled!" + ColorReset)
			} else {
				peerTotal += score
				enc.Encode(NetMsg{T: "score", Idx: 1, Delta: score, Total: peerTotal})
				if peerTotal >= WinningScore {
					enc.Encode(NetMsg{T: "game_over", Idx: 1})
					fmt.Println(ColorRed + "💀 Peer wins. Returning to menu." + ColorReset)
					return
				}
				if act.Bank {
					diceToRoll = 6
				} else {
					if len(act.Keep) == diceToRoll {
						diceToRoll = 6
						encHot <- NetMsg{T: "hot", Idx: 1}
					} else {
						diceToRoll -= len(act.Keep)
					}
					continue
				}
			}
		}
		currentIdx = 1 - currentIdx
	}
}

//MARK: Host Turn Loop
func hostTurnLoop(roll []int, diceToRoll int, encHot chan NetMsg) (turnScore int, diceRemaining int, turnEnded bool) {
	if calculateScore(roll) == 0 {
		fmt.Println(ColorRed + "Farkle! You scored 0 this turn." + ColorReset)
		encHot <- NetMsg{T: "farkle", Idx: 0}
		return 0, 6, true
	}
	promptText := ColorBlue + "Commands: 'keep X X...' (score & continue), 'bank X X...' (score & pass), or 'quit'" + ColorReset
	for {
		fmt.Println(promptText)
		kept, action := promptAction(roll)
		switch action {
		case "quit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		case "keep":
			score := calculateScore(kept)
			turnScore += score
			fmt.Printf(ColorGreen+"Scored %d (turn total %d)."+ColorReset+"\n", score, turnScore)
			if len(kept) == diceToRoll {
				fmt.Println(ColorYellow + "Hot dice!" + ColorReset)
				encHot <- NetMsg{T: "hot", Idx: 0}
				newRoll := rollDice(6)
				renderDice(newRoll)
				encHot <- NetMsg{T: "roll", Dice: newRoll, Idx: 0}
				returnScore, newRemaining, _ := hostTurnLoop(newRoll, 6, encHot)
				return turnScore + returnScore, newRemaining, false
			}
			returnScore := turnScore
			newDice := diceToRoll - len(kept)
			newRoll := rollDice(newDice)
			renderDice(newRoll)
			encHot <- NetMsg{T: "roll", Dice: newRoll, Idx: 0}
			returnScore2, newRemaining, _ := hostTurnLoop(newRoll, newDice, encHot)
			return returnScore + returnScore2, newRemaining, false
		case "bank":
			score := calculateScore(kept)
			turnScore += score
			fmt.Printf(ColorGreen+"Banking %d (turn total %d)."+ColorReset+"\n", score, turnScore)
			return turnScore, 6, true
		}
	}
}

//MARK: Peer Lobby
func JoinLobby(hostIP, lobbyID string) {
	addr := net.JoinHostPort(hostIP, "9313")
	fmt.Println("Dialling", addr, "with lobby ID", lobbyID, "…")

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		fmt.Println(ColorRed+"Connection failed:", err, ColorReset)
		return
	}
	defer conn.Close()

	enc := json.NewEncoder(conn)
	dec := json.NewDecoder(conn)

	enc.Encode(NetMsg{T: "hello", Name: "peer"})
	var welcome NetMsg
	if err := dec.Decode(&welcome); err != nil || welcome.T != "welcome" {
		fmt.Println(ColorRed + "Handshake failed." + ColorReset)
		return
	}
	WinningScore = welcome.Target
	fmt.Println(ColorGreen+"Connected! Target score:", WinningScore, ColorReset)

	var mu sync.Mutex
	var lastRoll []int

	for {
		var msg NetMsg
		if err := dec.Decode(&msg); err != nil {
			fmt.Println(ColorRed + "Connection lost." + ColorReset)
			return
		}

		switch msg.T {
		case "roll":
			renderDice(msg.Dice)
			lastRoll = msg.Dice
		case "your_turn":
			mu.Lock()
			turnLoopPeer(enc, lastRoll)
			mu.Unlock()
		case "farkle":
			if msg.Idx == 0 {
				fmt.Println(ColorRed + "Host Farkled." + ColorReset)
			} else {
				fmt.Println(ColorRed + "You Farkled." + ColorReset)
			}
		case "score":
			if msg.Idx == 0 {
				fmt.Printf(ColorYellow+"Host scored %d (total %d)"+ColorReset+"\n", msg.Delta, msg.Total)
			} else {
				fmt.Printf(ColorGreen+"You scored %d (total %d)"+ColorReset+"\n", msg.Delta, msg.Total)
			}
		case "game_over":
			if msg.Idx == 1 {
				fmt.Println(ColorGreen + "🏆 You win! Returning to menu." + ColorReset)
			} else {
				fmt.Println(ColorRed + "💀 Host wins. Returning to menu." + ColorReset)
			}
			return
		case "hot":
			if msg.Idx == 0 {
				fmt.Println(ColorYellow + "Host got hot dice!" + ColorReset)
			} else {
				fmt.Println(ColorYellow + "Hot dice! Rolling all 6 again..." + ColorReset)
			}
		case "banner":
			fmt.Printf("\n========================\n")
			fmt.Printf(" ROUND %d – First to %d\n", msg.Round, msg.Target)
			fmt.Printf("========================\n")
			fmt.Printf("Scoreboard → %sHost%s: %d | %sYou%s: %d\n", ColorYellow, ColorReset, msg.HostTotal, ColorGreen, ColorReset, msg.PeerTotal)
			fmt.Println("\n" + ColorRed + "Host turn:" + ColorReset)
		}
	}
}

func turnLoopPeer(enc *json.Encoder, roll []int) {
	promptText := ColorBlue + "Commands: 'keep X X...' (score & continue), 'bank X X...' (score & pass), or 'quit'" + ColorReset
	for {
		fmt.Println(promptText)
		kept, action := promptAction(roll)
		switch action {
		case "quit":
			fmt.Println("Goodbye!")
			os.Exit(0)
		case "keep":
			enc.Encode(NetMsg{T: "action", Keep: kept, Bank: false})
			return
		case "bank":
			enc.Encode(NetMsg{T: "action", Keep: kept, Bank: true})
			return
		}
	}
}

//------------------------------------------------------------
func generateLobbyID() string {
	buf := make([]byte, 5)
	crand.Read(buf)
	return strings.ToUpper(b32.EncodeToString(buf))
}