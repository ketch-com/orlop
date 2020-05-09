// Copyright (c) 2020 SwitchBit, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package orlop

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/switch-bit/orlop/log"
)

// HasCredentialsConfig denotes that an object provides credentials configuration
type HasCredentialsConfig interface {
	GetID() string
	GetUsername() string
	GetPassword() string
}

// CredentialsConfig provides credentials configuration
type CredentialsConfig struct {
	ID       string
	Username string
	Password string
}

// GetID returns the ID of the credentials in Vault
func (c CredentialsConfig) GetID() string {
	return c.ID
}

// GetUsername returns a static username
func (c CredentialsConfig) GetUsername() string {
	return c.Username
}

// GetPassword returns a static password
func (c CredentialsConfig) GetPassword() string {
	return c.Password
}

// Credentials provides username/password information
type Credentials struct {
	Username string
	Password string
}

// GetCredentials retrieves credentials
func GetCredentials(cfg HasCredentialsConfig, vault HasVaultConfig) (*Credentials, error) {
	l := log.WithFields(logrus.Fields{
		"credentials.id": cfg.GetID(),
		"vault.enabled":  vault != nil && vault.GetEnabled(),
	})
	l.Debug("loading credentials")

	if len(cfg.GetUsername()) > 0 && len(cfg.GetPassword()) > 0 {
		l.Trace("loaded from inline settings")

		return &Credentials{
			Username: cfg.GetUsername(),
			Password: cfg.GetPassword(),
		}, nil
	}

	if len(cfg.GetID()) == 0 || vault == nil || !vault.GetEnabled() {
		l.Error("no credentials specified")
		return nil, fmt.Errorf("credentials: no credentials specified")
	}

	client, err := NewVault(vault)
	if err != nil {
		l.WithError(err).Error("could not connect to Vault")
		return nil, err
	}

	s, err := client.Read(cfg.GetID())
	if err != nil {
		l.WithError(err).Error("credentials not found")
		return nil, err
	}

	if s == nil {
		return nil, fmt.Errorf("credentials: could not load credentials from %s", cfg.GetID())
	} else if s.Data["username"] == nil {
		l.Tracef("username not found")
		return nil, fmt.Errorf("credentials: could not load credentials from %s", cfg.GetID())
	} else if s.Data["password"] == nil {
		l.Tracef("password not found")
		return nil, fmt.Errorf("credentials: could not load credentials from %s", cfg.GetID())
	}

	creds := &Credentials{}
	if u, ok := s.Data["username"].(string); ok {
		creds.Username = u
	}
	if p, ok := s.Data["password"].(string); ok {
		creds.Password = p
	}

	return creds, nil
}

// AddCredentials adds credentials configuration parameters
func AddCredentials(flags *pflag.FlagSet, prefix ...string) {
	p := MakeCommandKeyPrefix(prefix)
	flags.String(p("id"), "", "Vault path to retrieve credentials")
	flags.String(p("username"), "", "static username")
	flags.String(p("password"), "", "static password")
}
