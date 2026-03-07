package engine

type BetStrategy interface {
	Type() BetType
	Validate(b Bet) error
	IsWin(d DiceResult, b Bet) bool
	Odds(b Bet, d DiceResult) int64
}
