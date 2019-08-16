package loggregator

import (
	"code.cloudfoundry.org/go-loggregator"
	"github.com/pkg/errors"
)

func NewIngressClient(ca, cert, key, addr string) (*loggregator.IngressClient, error) {
	tlsConfig, err := loggregator.NewIngressTLSConfig(ca, cert, key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create forwarder")
	}

	return loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr(addr),
	)
}
