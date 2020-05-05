package orlop

import (
	"crypto/tls"
	"fmt"
)

func newClientAuthValue() *clientAuthValue {
	return &clientAuthValue{
		v: tls.NoClientCert,
	}
}

type clientAuthValue struct {
	v tls.ClientAuthType
}

func (c clientAuthValue) String() string {
	return fmt.Sprintf("%d", c.v)
}

func (c *clientAuthValue) Set(s string) error {
	switch s {
	case "request":
		c.v = tls.RequestClientCert

	case "require":
		c.v = tls.RequireAnyClientCert

	case "verify":
		c.v = tls.VerifyClientCertIfGiven

	case "require-and-verify":
		c.v = tls.RequireAndVerifyClientCert

	case "none":
		c.v = tls.NoClientCert

	default:
		return fmt.Errorf("config: '%s' is not a valid ClientAuth value", s)
	}

	return nil
}

func (c clientAuthValue) Type() string {
	return "ClientAuthType"
}
