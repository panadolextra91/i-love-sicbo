package engine

import "sync"

var settleMu sync.Mutex

func SettleRoundAtomic(roundID string, bets []Bet, roller DiceRoller, reg *StrategyRegistry, wallet WalletStore, roundRepo RoundRepo) (RoundSettlement, error) {
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

	if err := wallet.ApplyGrossBatch(playerGross); err != nil {
		return RoundSettlement{}, err
	}

	if err := roundRepo.MarkSettled(roundID, dice, results); err != nil {
		return RoundSettlement{}, err
	}

	return RoundSettlement{
		RoundID:     roundID,
		Dice:        dice,
		PlayerGross: playerGross,
		Details:     results,
	}, nil
}
