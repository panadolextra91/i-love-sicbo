package network

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/mdns"
)

type MDNSServer struct {
	instanceName string
	server       *mdns.Server
}

func shortUUID() string {
	u := uuid.NewString()
	if len(u) < 8 {
		return u
	}
	return u[:8]
}

func StartMDNSAdvertise(cfg Config) (*MDNSServer, error) {
	instance := fmt.Sprintf("Cachon-Casino-%s", shortUUID())
	service, err := mdns.NewMDNSService(instance, cfg.ServiceName, "", "", cfg.WSPort, nil, []string{"instance=" + instance})
	if err != nil {
		return nil, err
	}

	srv, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return nil, err
	}

	return &MDNSServer{instanceName: instance, server: srv}, nil
}

func (m *MDNSServer) InstanceName() string {
	return m.instanceName
}

func (m *MDNSServer) Stop() {
	if m != nil && m.server != nil {
		_ = m.server.Shutdown()
	}
}
