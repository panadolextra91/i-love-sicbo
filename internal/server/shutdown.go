package server

import (
	"time"

	"cachon-casino/internal/network"
)

func BroadcastShutdown(deps Deps, reason string, drain time.Duration) {
	env, _ := network.NewEnvelope(network.MsgServerShutdown, deps.SessionID, deps.Seq.Next(), network.ShutdownPayload{
		Reason:     reason,
		ShutdownAt: time.Now().Add(drain).UnixMilli(),
	})
	deps.Hub.Broadcast <- env
}
