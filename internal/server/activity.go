package server

import (
	"fmt"
	"math/rand"
	"time"

	"cachon-casino/internal/engine"
)

type ActivityEngine struct {
	r *rand.Rand
}

func NewActivityEngine() *ActivityEngine {
	return &ActivityEngine{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (a *ActivityEngine) Build(settlement engine.RoundSettlement) []string {
	phrases := []string{
		"vừa húp chip ngọt như mía lùi",
		"ăn kèo xong cười khẩy cả bàn",
		"nổ kèo làm sòng im re 3 giây",
	}
	out := make([]string, 0, len(settlement.PlayerGross))
	for pid, gross := range settlement.PlayerGross {
		if gross <= 0 {
			continue
		}
		out = append(out, fmt.Sprintf("%s %s: +%d chip", pid, phrases[a.r.Intn(len(phrases))], gross))
	}
	return out
}
