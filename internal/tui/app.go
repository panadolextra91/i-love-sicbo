package tui

import (
	"context"
	"time"

	"cachon-casino/internal/engine"
	"cachon-casino/internal/network"

	tea "github.com/charmbracelet/bubbletea"
)

type ModelState string

const (
	StateJoining  ModelState = "joining"
	StateBetting  ModelState = "betting"
	StateRolling  ModelState = "rolling"
	StateSettled  ModelState = "settled"
	StateShutdown ModelState = "shutdown"
)

type Model struct {
	State       ModelState
	Config      network.Config
	SessionID   string
	PlayerID    string
	PlayerName  string
	Chips       int64
	RoundID     string
	Countdown   int
	BettingOpen bool

	SelectedBetType engine.BetType
	BetStake        int64
	TargetValue     int

	Dice        [3]int
	History     []string // 🌸 (Tài) / 🍀 (Xỉu)
	LastResults []engine.PayoutResult

	Logs []string
	Err  error

	// Conn and channel for WebSocket
	Inbound  chan network.Envelope
	Outbound chan network.Envelope
	Ctx      context.Context
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.EnterAltScreen,
		m.waitForMsg(),
	)
}

func (m Model) waitForMsg() tea.Cmd {
	return func() tea.Msg {
		return <-m.Inbound
	}
}

type TickMsg time.Time

type ResetRollingMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
		return m.handleKeys(msg)

	case network.Envelope:
		return m.handleEnvelope(msg)

	case TickMsg:
		if m.State == StateRolling {
			// Continue animation tick
			return m, tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg {
				return TickMsg(t)
			})
		}

	case ResetRollingMsg:
		m.State = StateSettled
		// Update history
		sum := m.Dice[0] + m.Dice[1] + m.Dice[2]
		if sum >= 11 {
			m.History = append(m.History, "🌸")
		} else {
			m.History = append(m.History, "🍀")
		}
		if len(m.History) > 30 {
			m.History = m.History[1:]
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleEnvelope(env network.Envelope) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd = m.waitForMsg()

	switch env.Type {
	case network.MsgJoinAck:
		var p network.JoinAckPayload
		_ = network.DecodePayloadTo(env, &p)
		m.PlayerID = p.PlayerID
		m.PlayerName = p.Name
		m.Chips = p.Chips
		m.RoundID = p.RoundID
		m.State = StateBetting

	case network.MsgCountdownTick:
		var p network.CountdownPayload
		_ = network.DecodePayloadTo(env, &p)
		m.RoundID = p.RoundID
		m.Countdown = p.SecondsLeft
		m.BettingOpen = p.BettingOpen

	case network.MsgRoundResult:
		var p network.RoundResultPayload
		_ = network.DecodePayloadTo(env, &p)
		m.Dice = p.Dice
		m.State = StateRolling
		// Start rolling animation
		cmd = tea.Batch(
			cmd,
			tea.Tick(150*time.Millisecond, func(t time.Time) tea.Msg { return TickMsg(t) }),
			func() tea.Msg {
				time.Sleep(2 * time.Second)
				return ResetRollingMsg{}
			},
		)

	case network.MsgBetAccepted:
		// Optional: Flash something green
	}

	return m, cmd
}

func (m Model) View() string {
	if m.State == StateJoining {
		return "Joining the casino... 🌸"
	}
	return m.renderDashboard()
}
