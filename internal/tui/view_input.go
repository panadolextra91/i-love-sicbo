package tui

import (
	"cachon-casino/internal/engine"
	"cachon-casino/internal/network"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "1":
		if m.SelectedBetType == engine.BetBig {
			m.SelectedBetType = engine.BetSmall
		} else {
			m.SelectedBetType = engine.BetBig
		}
	case "2":
		if m.SelectedBetType == engine.BetOdd {
			m.SelectedBetType = engine.BetEven
		} else {
			m.SelectedBetType = engine.BetOdd
		}
	case "3":
		m.SelectedBetType = engine.BetAnyTriple
	case "4":
		m.SelectedBetType = engine.BetExactTriple
	case "5":
		m.SelectedBetType = engine.BetExactDouble
	case "6":
		m.SelectedBetType = engine.BetTwoNumberCombo
	case "7":
		m.SelectedBetType = engine.BetExactTotal
	case "8":
		m.SelectedBetType = engine.BetSingleNumber
	case "tab":
		m.TargetValue++
		if m.TargetValue > 17 {
			m.TargetValue = 1
		}
	case "+":
		m.BetStake += 100
	case "-":
		if m.BetStake >= 100 {
			m.BetStake -= 100
		}
	case "enter":
		if m.BettingOpen && m.BetStake > 0 {
			return m, m.placeBet()
		}
	}
	return m, nil
}

func (m Model) placeBet() tea.Cmd {
	return func() tea.Msg {
		bet := network.PlaceBetPayload{
			RoundID:     m.RoundID,
			BetType:     string(m.SelectedBetType),
			Stake:       m.BetStake,
			TargetValue: m.TargetValue,
		}
		env, _ := network.NewEnvelope(network.MsgPlaceBet, m.SessionID, 0, bet)
		m.Outbound <- env
		return nil
	}
}
