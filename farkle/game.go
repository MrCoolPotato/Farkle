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

const (
    ColorReset  = "\033[0m"
    ColorRed    = "\033[31m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorCyan   = "\033[36m"
    ColorBlue   = "\033[34m"
)

var WinningScore = 1000
var aiDelay = 2 * time.Second

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

// MARK: Main game loop
func PlayGame() {
    rand.Seed(time.Now().UnixNano())

    player := Player{Name: "You"}
    enemy := Player{Name: "Enemy"}

    round := 1
    for {
        fmt.Printf("\n========================\n")
        fmt.Printf(" ROUND %d – First to %d\n", round, WinningScore)
        fmt.Printf("========================\n")
        fmt.Printf("Scoreboard → %sYou%s: %d | %sEnemy%s: %d\n",
            ColorGreen, ColorReset, player.Total,
            ColorRed, ColorReset, enemy.Total)

        fmt.Println("\n" + ColorGreen + "Your turn:" + ColorReset)
        fmt.Println(ColorYellow + "First roll will happen automatically; then choose dice to keep or bank." + ColorReset)
        playerPoints := playerTurn()
        player.Total += playerPoints
        fmt.Printf("You banked %d points. New total: %d\n", playerPoints, player.Total)
        if player.Total >= WinningScore {
            fmt.Println("\n" + ColorGreen + "🏆  You win the game! 🏆" + ColorReset)
            return
        }

        fmt.Println("\n" + ColorRed + "Enemy turn:" + ColorReset)
        enemyPoints := enemyTurn()
        enemy.Total += enemyPoints
        fmt.Printf("Enemy banked %d points. New total: %d\n", enemyPoints, enemy.Total)
        if enemy.Total >= WinningScore {
            fmt.Println("\n" + ColorRed + "💀  Enemy wins the game. Better luck next time!" + ColorReset)
            return
        }

        round++
    }
}

// MARK: Dice rolling
func rollDice(n int) []int {
    dice := make([]int, n)
    for i := 0; i < n; i++ {
        dice[i] = rand.Intn(6) + 1
    }
    return dice
}

// MARK: Dice rendering
func renderDice(dice []int) {
    for _, d := range dice {
        fmt.Printf("%s[%s %d]%s ", ColorCyan, dieFaces[d], d, ColorReset)
    }
    fmt.Println()
}

// MARK: Score calculation
func calculateScore(dice []int) int {
    counts := make(map[int]int)
    for _, d := range dice {
        counts[d]++
    }

    if len(dice) == 6 {
        full := true
        for i := 1; i <= 6; i++ {
            if counts[i] != 1 {
                full = false
                break
            }
        }
        if full {
            return 1500
        }
    }

    if len(dice) == 5 {
        isStraight15 := true
        for i := 1; i <= 5; i++ {
            if counts[i] != 1 {
                isStraight15 = false
                break
            }
        }
        if isStraight15 {
            return 500
        }
        isStraight26 := true
        for i := 2; i <= 6; i++ {
            if counts[i] != 1 {
                isStraight26 = false
                break
            }
        }
        if isStraight26 {
            return 750
        }
    }

    score := 0
    for val, cnt := range counts {
        if cnt >= 3 {
            base := 0
            if val == 1 {
                base = 1000
            } else {
                base = val * 100
            }
            mult := 1 << (cnt - 3)
            score += base * mult
            cnt = 0
        }
        if val == 1 {
            score += cnt * 100
        } else if val == 5 {
            score += cnt * 50
        }
    }
    return score
}

// MARK: AI scoring dice selection
func aiSelectScoringDice(roll []int) (kept []int) {
    counts := make(map[int]int)
    for _, d := range roll {
        counts[d]++
    }

    full := true
    for i := 1; i <= 6; i++ {
        if counts[i] != 1 {
            full = false
            break
        }
    }
    if full {
        return roll
    }

    if len(roll) >= 5 {
        part15 := true
        for i := 1; i <= 5; i++ {
            if counts[i] == 0 {
                part15 = false
                break
            }
        }
        if part15 {
            for _, d := range roll {
                if d != 6 {
                    kept = append(kept, d)
                }
            }
            return
        }
        part26 := true
        for i := 2; i <= 6; i++ {
            if counts[i] == 0 {
                part26 = false
                break
            }
        }
        if part26 {
            for _, d := range roll {
                if d != 1 {
                    kept = append(kept, d)
                }
            }
            return
        }
    }

    for val, cnt := range counts {
        if cnt >= 3 {
            for i := 0; i < cnt; i++ {
                kept = append(kept, val)
            }
        } else if val == 1 || val == 5 {
            for i := 0; i < cnt; i++ {
                kept = append(kept, val)
            }
        }
    }
    return
}

// MARK: Player action prompt
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
            if calculateScore(vals) == 0 {
                fmt.Println("Selected dice do not form a scoring combination.")
                return nil, false
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

        fmt.Println(ColorBlue + "Commands: 'keep X X...' (score & continue), 'bank X X...' (score & pass), or 'quit'" + ColorReset)
    }
}

// MARK: Player turn logic
func playerTurn() int {
    diceToRoll := 6
    turnScore := 0

outer:
    for {
        fmt.Printf("-- Rolling %d dice --\n", diceToRoll)
        roll := rollDice(diceToRoll)
        renderDice(roll)

        if calculateScore(roll) == 0 {
            fmt.Println(ColorRed + "Farkle! You lose all unbanked points for this turn." + ColorReset)
            return 0
        }

        fmt.Println(ColorBlue + "Commands: 'keep X X...' to score & CONTINUE, 'bank X X...' to score & PASS, or 'quit'" + ColorReset)

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
                fmt.Printf(ColorGreen+"Scored %d (turn total %d). Continuing..."+ColorReset+"\n", score, turnScore)

                diceToRoll -= len(kept)
                if diceToRoll == 0 {
                    fmt.Println(ColorYellow + "Hot dice! All dice scored, rolling 6 fresh dice." + ColorReset)
                    diceToRoll = 6
                }
                break inner
            case "bank":
                score := calculateScore(kept)
                turnScore += score
                fmt.Printf(ColorGreen+"Banking %d points (turn total %d)."+ColorReset+"\n", score, turnScore)
                return turnScore
            }
        }
        continue outer
    }
}

// MARK: Enemy turn logic
func enemyTurn() int {
    diceToRoll := 6
    turnScore := 0

    for {
        time.Sleep(aiDelay)
        fmt.Printf("-- Enemy rolling %d dice --\n", diceToRoll)
        time.Sleep(aiDelay)

        roll := rollDice(diceToRoll)
        renderDice(roll)

        if calculateScore(roll) == 0 {
            fmt.Println(ColorRed + "Enemy Farkled and scores 0." + ColorReset)
            return 0
        }

        kept := aiSelectScoringDice(roll)
        score := calculateScore(kept)
        turnScore += score
        fmt.Printf("Enemy keeps %v gaining %d (turn total %d).\n", kept, score, turnScore)

        diceToRoll -= len(kept)
        if diceToRoll == 0 {
            fmt.Println(ColorYellow + "Enemy got hot dice and will roll all 6 again!" + ColorReset)
            diceToRoll = 6
        }

        if turnScore >= 1000 || diceToRoll <= 2 || len(kept) >= 5 {
            time.Sleep(aiDelay)
            fmt.Println(ColorBlue + "Enemy decides to bank." + ColorReset)
            return turnScore
        }
    }
}