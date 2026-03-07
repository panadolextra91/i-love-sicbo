package network

import "time"

type Config struct {
	ServiceName      string
	WSPort           int
	ClientSendBuffer int
	PingInterval     time.Duration
	PongWait         time.Duration
	AppPingInterval  time.Duration
	HeartbeatTimeout time.Duration
	LatencyBuffer    time.Duration
}

func DefaultConfig() Config {
	return Config{
		ServiceName:      "_cachon._tcp",
		WSPort:           8080,
		ClientSendBuffer: 256,
		PingInterval:     5 * time.Second,
		PongWait:         10 * time.Second,
		AppPingInterval:  5 * time.Second,
		HeartbeatTimeout: 10 * time.Second,
		LatencyBuffer:    500 * time.Millisecond,
	}
}
