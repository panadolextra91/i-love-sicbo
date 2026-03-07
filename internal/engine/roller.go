package engine

import (
	"crypto/rand"
	"math/big"
)

type DiceRoller interface {
	Roll3() DiceResult
}

type CryptoDiceRoller struct{}

func (r CryptoDiceRoller) Roll3() DiceResult {
	return DiceResult{secureDie(), secureDie(), secureDie()}
}

func secureDie() int {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		panic(err)
	}
	return int(n.Int64()) + 1
}

type FixedDiceRoller struct {
	result DiceResult
}

func NewFixedDiceRoller(d DiceResult) DiceRoller {
	return FixedDiceRoller{result: d}
}

func (r FixedDiceRoller) Roll3() DiceResult {
	return r.result
}
