package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderDashboard() string {
	// 1. Header
	logo := `
  ___   _       ___ __     __ _____   ____ ___  ____     ____   ___  
 |_ _| | |     / _ \\ \   / /| ____| / ___|_ _|/ ___|   | __ ) / _ \ 
  | |  | |    | | | |\ \ / / |  _|   \___ \| | | |       |  _ \| | | |
  | |  | |___ | |_| | \ V /  | |___   ___) | | | |___    | |_) | |_| |
 |___| |_____| \___/   \_/   |_____| |____/___| \____|   |____/ \___/ 
`
	header := lipgloss.JoinVertical(lipgloss.Center,
		TitleStyle.Render(logo),
		SloganStyle.Render("(ĐÁNH TÀI RA TÀI, ĐÁNH XỈU RA XỈU)"),
	)

	// 2. Game Info (Round, Chips, Timer)
	info := lipgloss.JoinHorizontal(lipgloss.Center,
		BorderPink.Render(fmt.Sprintf("Ván: #%s", m.RoundID)),
		BorderPink.Render(fmt.Sprintf("Ví: %s 💰", ChipStyle.Render(fmt.Sprintf("%d", m.Chips)))),
		BorderPink.Render(fmt.Sprintf("Thời gian: %s", TimerStyle.Render(fmt.Sprintf("%ds", m.Countdown)))),
	)

	// 3. Middle Section (Dice & History)
	historyStr := "Lịch sử: " + strings.Join(m.History, " ")
	middle := lipgloss.JoinVertical(lipgloss.Center,
		m.renderDice(),
		lipgloss.NewStyle().PaddingTop(1).Render(historyStr),
		lipgloss.NewStyle().PaddingTop(1).Render(m.LastOutcomeText),
		func() string {
			if m.State == StateSettled {
				return lipgloss.NewStyle().PaddingTop(1).Render("Nhấn F để tiếp tục vòng mới")
			}
			return ""
		}(),
	)

	// 4. Betting Section
	betting := m.renderBettingArea()

	// 5. Combine All
	main := lipgloss.JoinVertical(lipgloss.Center,
		header,
		info,
		middle,
		betting,
	)

	// Center everything in the terminal
	return main
}

func (m Model) renderBettingArea() string {
	title := TitleStyle.Render("--- ĐẶT CƯỢC ---")

	options := []string{
		"1: Tài/Xỉu",
		"2: Chẵn/Lẻ",
		"3: Bộ ba b.kỳ",
		"4: Bộ ba CX",
		"5: Đôi CX",
		"6: Cặp số",
		"7: Tổng CX",
		"8: Một số",
	}

	betRow := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Padding(1).Render(strings.Join(options, " | ")),
	)

	status := fmt.Sprintf("Đang chọn: %s | Tiền: %d | Số: %d",
		BetTypeStyle.Render(string(m.SelectedBetType)),
		m.BetStake,
		m.TargetValue,
	)

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		betRow,
		lipgloss.NewStyle().Foreground(SoftPink).Render(status),
		lipgloss.NewStyle().PaddingTop(1).Foreground(Gray).Render("Phím: [1-8] Chọn cửa, [+/-] Tiền, [Tab] Đổi số, [Enter] Cược, [Q] Thoát"),
	)
}
