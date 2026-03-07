package engine

import "sync"

var settleMu sync.Mutex

func SettleRoundAtomic(roundID string, roundNo int64, startedAt, settledAt int64, bets []Bet, roller DiceRoller, reg *StrategyRegistry, store SettlementStore) (RoundSettlement, error) {
	settleMu.Lock()
	defer settleMu.Unlock()

	dice := roller.Roll3()
	results := make([]PayoutResult, 0, len(bets))
	playerGross := make(map[string]int64)

	for _, b := range bets {
		gross := CalculatePayout(reg, dice[:], b)
		results = append(results, PayoutResult{Bet: b, Win: gross > 0, GrossPayout: gross})
		playerGross[b.PlayerID] += gross
	}

	if err := store.SettleRound(roundID, roundNo, startedAt, settledAt, dice, results); err != nil {
		return RoundSettlement{}, err
	}

	return RoundSettlement{
		RoundID:     roundID,
		Dice:        dice,
		PlayerGross: playerGross,
		Details:     results,
	}, nil
}
