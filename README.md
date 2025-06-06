![Farkle Logo](farkle_raw.png)

# Farkle (CLI Edition)

Welcome to **Farkle** – a fully-playable terminal implementation of the dice game.
The project offers both a single-player mode (human vs. AI) and a peer-to-peer multiplayer mode with automatic lobby IDs and optional UPnP port-mapping for hassle-free Internet play.

---

## Table of Contents

1. [Features](#features)
2. [Gameplay Quick-Start](#gameplay-quick-start)
3. [Installation & Build](#installation--build)
4. [Running the Game](#running-the-game)
5. [Multiplayer Details](#multiplayer-details)
6. [Scoring Reference](#scoring-reference)
7. [Project Structure](#project-structure)

---

## Features

* **Game rules**: 1s & 5s, triples, 4-6-of-a-kind multipliers, 1-5 & 2-6 straights, full straight, hot-dice and farkle busts.
* **Enemy AI**: basic risk heuristic; banks intelligently and recognises all scoring combos.
* **Colourful TUI**: distinct colours for banners, dice, prompts, peer rolls, hot-dice & farkles.
* **Keep / Bank commands** exactly like they sound *Score & Continue* / *Score & Pass*.
* **Peer-to-Peer multiplayer**:
  * Auto-generated Lobby ID encodes host’s public IPv4 + port.
  * Optional UPnP port-mapping (TCP 9313) – no manual router config in most home networks.
  * Live ping keep-alive to detect disconnects.
* **Configurable winning score** (`play 15000` → first to 15 000).

---

## Gameplay Quick-Start

| Command                     | Result                                           |
| --------------------------- | ------------------------------------------------ |
| `play`                      | Solo game to 1 000 points.                       |
| `play 10000`                | Solo to 10 000.                                  |
| `play --mp --create`        | Host a lobby, prints Lobby ID (e.g. `B4Q5FPHG`). |
| `play --mp --join=B4Q5FPHG` | Join that lobby – IP+port decoded automatically. |
| `keep 1 5 5`                | Score those dice & continue.                     |
| `bank 1 1 1`                | Score & pass turn.                               |
| `quit` / `exit`             | Leave at any prompt.                             |

---

## Installation & Build

Refer to the packages section or continue down.

```bash
# Clone and install deps
$ git clone
$ cd farkle-go
$ go mod tidy   # fetch deps

# Run directly
$ go run .

# Or build binary
$ go build -o farkle
```

Go latest+ recommended.

---

## Running the Game

```bash
# Single-player (default score 1 000)
$ ./farkle play

# Single-player to 20 000
$ ./farkle play 20000

# Multiplayer host
$ ./farkle play --mp --create
#  UPnP mapped external port 9313
#  Lobby created. Share ID: B4Q5FPHG

# Peer joins (no extra flags needed)
$ ./farkle play --mp --join=B4Q5FPHG
```

If UPnP fails you will see a yellow notice – forward TCP 9313 manually.

---

## Multiplayer Details

* **Lobby ID** – Base-32 encodes host IPv4 (4 B) + external port (2 B) + random byte.
* **Control channel** – plain TCP (9313). Host authoritative.
* **Ping** – 10 s heartbeat, 30 s timeout.
* **Security** – plaintext.

---

## Scoring Reference

| Combination                  | Points                                            |
| ---------------------------- | ------------------------------------------------- |
| Single **1**                 | 100                                               |
| Single **5**                 | 50                                                |
| Three 1s                     | 1 000                                             |
| Three 2s / 3s / 4s / 5s / 6s | 200 – 600                                         |
| 4-of-a-kind                  | double triple                                     |
| 5-of-a-kind                  | ×4 triple                                         |
| 6-of-a-kind                  | ×8 triple                                         |
| Straight 1-5                 | 500                                               |
| Straight 2-6                 | 750                                               |
| Straight 1-6                 | 1 500                                             |
| **Farkle**                   | Roll with *zero* scoring combos – lose turn score |
| **Hot Dice**                 | All dice score – roll fresh 6 & continue          |

---

## Project Structure

```
Farkle/
├─farkle/
    ├─ game.go        # single-player logic
    └─ game_mp.go     # multi-player logic
├─ main.go        # CLI menu & flag parsing
├─ go.mod / sum   # module file
└─ README.md
```

---