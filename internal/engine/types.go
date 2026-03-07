package engine

type BetType string

const (
	BetBig            BetType = "big"
	BetSmall          BetType = "small"
	BetEven           BetType = "even"
	BetOdd            BetType = "odd"
	BetAnyTriple      BetType = "any_triple"
	BetExactTriple    BetType = "exact_triple"
	BetExactDouble    BetType = "exact_double"
	BetTwoNumberCombo BetType = "two_number_combo"
	BetExactTotal     BetType = "exact_total"
	BetSingleNumber   BetType = "single_number"
)

type Bet struct {
	PlayerID    string
	Type        BetType
	Stake       int64
	TargetValue int
}

type DiceResult [3]int

type PayoutResult struct {
	Bet         Bet
	Win         bool
	GrossPayout int64
}

type RoundSettlement struct {
	RoundID     string
	Dice        DiceResult
	PlayerGross map[string]int64
	Details     []PayoutResult
}

type WalletStore interface {
	ApplyGrossBatch(playerGross map[string]int64) error
}

type RoundRepo interface {
	MarkSettled(roundID string, dice DiceResult, results []PayoutResult) error
}

type BetStore interface {
	Append(bet Bet) error
}
