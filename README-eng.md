# 🌸 I LOVE SIC BO 🌸

> **Languages:** [English](README-eng.md) | [Tiếng Việt](README.md)

```text
  ___   _       ___ __     __ _____   ____ ___  ____     ____   ___  
 |_ _| | |     / _ \\ \   / /| ____| / ___|_ _|/ ___|   | __ ) / _ \ 
  | |  | |    | | | |\ \ / / |  _|   \___ \| | | |       |  _ \| | | |
  | |  | |___ | |_| | \ V /  | |___   ___) | | | |___    | |_) | |_| |
 |___| |_____| \___/   \_/   |_____| |____/___| \____|   |____/ \___/ 
```

Welcome to **I LOVE SIC BO** — a "premium" Sic Bo (Tai Xiu) casino built with a Terminal User Interface (TUI), featuring a stunning hot pink aesthetic! 🎀

This is a complete Go project that allows you to host a casino server over LAN and invite friends to join the gamble via a beautiful ASCII Art interface.

## ✨ Key Features

- **🎲 Pink TUI Interface**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lipgloss](https://github.com/charmbracelet/lipgloss) — vibrant pink colors with artistic ASCII fonts.
- **📡 Auto-Discovery**: Utilizes mDNS so clients can find the server on a LAN automatically without typing IP addresses.
- **🎰 8 Betting Types**: Big/Small, Even/Odd, Triple, Double, Pair, Exact Total... full features like a professional casino.
- **📈 History Board**: Track match history with 🌸 (Big) and 🍀 (Small) icons.
- **💾 Persistence**: Balances and transaction history are stored in SQLite (`casino.db`).
- **⚡ Animations**: Smooth dice rolling animations before revealing results.

## 🛠 Prerequisites

- **Language**: Go version 1.21+ (or latest).
- **OS**: macOS, Linux, or Windows (Recommended to use a terminal with 256 colors/TrueColor support).

## 🚀 Installation & Usage

### 1. Clone & Install Dependencies

```bash
git clone https://github.com/panadolextra91/i-love-sicbo.git
cd i-love-sicbo
go mod tidy
```

### 2. Run the Server (Dealer)

The host needs to run:

```bash
go run cmd/server/main.go
```

> **Note**: The first run will automatically create `casino.db`.

### 3. Run the Client (Player)

All players (including the host) run:

```bash
go run cmd/client/main.go
```

The client will automatically find the server. If mDNS fails, it fallbacks to `localhost`.

## 🎮 Controls & Shortcuts

Control the game "Hacker style" using hotkeys:

| Key | Action |
|-----|--------|
| `[1]` | Bet **Big/Small** (Press repeatedly to toggle) |
| `[2]` | Bet **Even/Odd** |
| `[3]` | **Any Triple** |
| `[4]` | **Specific Triple** (111 - 666) |
| `[5]` | **Specific Double** |
| `[6]` | **Different Pair** |
| `[7]` | **Exact Total** (4 - 17) |
| `[8]` | **Specific Number** (Guess occurrences 1-3) |
| `[Tab]` | Switch selection value (e.g., toggle Big to Small, or change the total number) |
| `[+]` / `[-]` | Increase/Decrease bet amount (Step: 100) |
| `[Enter]` | **PLACE BET** |
| `[Q]` / `[Ctrl+C]` | Quit |

## 🏗 Tech Stack

- **Core**: Go (Golang)
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Networking**: [Gorilla WebSocket](https://github.com/gorilla/websocket) & [Hashicorp mDNS](https://github.com/hashicorp/mdns)
- **Database**: SQLite3

---

## License

MIT License

Copyright (c) 2026 I LOVE SIC BO contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

Enjoy your time at the pinkest casino in the terminal! 🐹💖🌸
