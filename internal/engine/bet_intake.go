package engine

func AcceptBet(bet Bet, reg *StrategyRegistry, store BetStore) error {
	strategy := reg.Get(bet.Type)
	if strategy == nil {
		return ErrUnknownBetType
	}

	if err := strategy.Validate(bet); err != nil {
		return err
	}

	return store.Append(bet)
}
