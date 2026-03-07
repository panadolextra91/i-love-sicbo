package network

import (
	"context"
	"net"

	"github.com/hashicorp/mdns"
)

type Endpoint struct {
	Instance string
	Host     string
	Port     int
}

func DiscoverEndpoint(ctx context.Context, serviceName string) (Endpoint, error) {
	entries := make(chan *mdns.ServiceEntry, 16)
	errCh := make(chan error, 1)

	go func() {
		errCh <- mdns.Lookup(serviceName, entries)
		close(entries)
	}()

	for {
		select {
		case <-ctx.Done():
			return Endpoint{}, ctx.Err()
		case err := <-errCh:
			if err != nil {
				return Endpoint{}, err
			}
		case e, ok := <-entries:
			if !ok {
				return Endpoint{}, context.DeadlineExceeded
			}
			if e == nil {
				continue
			}
			host := ""
			if e.AddrV4 != nil && !e.AddrV4.Equal(net.IP{}) {
				host = e.AddrV4.String()
			} else if e.Addr != nil {
				host = e.Addr.String()
			}
			if host == "" {
				continue
			}
			return Endpoint{Instance: e.Name, Host: host, Port: e.Port}, nil
		}
	}
}
