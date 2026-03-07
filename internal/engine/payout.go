package engine

func CalculatePayout(reg *StrategyRegistry, dice []int, bet Bet) int64 {
	if len(dice) != 3 {
		return 0
	}

	strategy := reg.Get(bet.Type)
	if strategy == nil {
		return 0
	}

	d := DiceResult{dice[0], dice[1], dice[2]}
	if strategy.IsWin(d, bet) {
		odds := strategy.Odds(bet, d)
		return bet.Stake * (odds + 1)
	}
	return 0
}
