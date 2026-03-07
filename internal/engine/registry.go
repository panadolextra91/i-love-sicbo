package engine

type StrategyRegistry struct {
	items map[BetType]BetStrategy
}

func NewStrategyRegistry(strategies ...BetStrategy) *StrategyRegistry {
	r := &StrategyRegistry{items: make(map[BetType]BetStrategy, len(strategies))}
	for _, s := range strategies {
		r.Register(s)
	}
	return r
}

func (r *StrategyRegistry) Register(s BetStrategy) {
	r.items[s.Type()] = s
}

func (r *StrategyRegistry) Get(t BetType) BetStrategy {
	return r.items[t]
}
