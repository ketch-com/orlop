// Copyright (c) 2020 Ketch Kloud, Inc.
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
	"context"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/v2/errors"
	"go.ketch.com/lib/orlop/v2/log"
)

// CredentialsConfig provides credentials configuration
type CredentialsConfig struct {
	ID       string
	Username string
	Password string
}

// Credentials provides username/password information
type Credentials struct {
	Username string
	Password string
}

// GetCredentials retrieves credentials
func GetCredentials(ctx context.Context, cfg CredentialsConfig, vault VaultConfig) (*Credentials, error) {
	ctx, span := tracer.Start(ctx, "GetCredentials")
	defer span.End()

	l := log.WithFields(logrus.Fields{
		"credentials.id": cfg.ID,
		"vault.enabled":  vault.GetEnabled(),
	})
	l.Trace("loading credentials")

	if len(cfg.Username) > 0 && len(cfg.Password) > 0 {
		l.Trace("loaded from inline settings")

		return &Credentials{
			Username: cfg.Username,
			Password: cfg.Password,
		}, nil
	}

	if len(cfg.ID) == 0 || !vault.GetEnabled() {
		err := errors.New("credentials: no credentials specified")
		span.RecordError(err)
		return nil, err
	}

	client, err := NewVault(ctx, vault)
	if err != nil {
		err := errors.Wrap(err, "credentials: could not connect to Vault")
		span.RecordError(err)
		return nil, err
	}

	s, err := client.Read(ctx, cfg.ID)
	if err != nil {
		err = errors.Wrap(err, "credentials: not found")
		span.RecordError(err)
		return nil, err
	}

	if s == nil {
		err := errors.Errorf("credentials: could not load credentials from %s", cfg.ID)
		span.RecordError(err)
		return nil, err
	} else if s.Data["username"] == nil {
		err := errors.Errorf("credentials: could not load credentials from %s", cfg.ID)
		span.RecordError(err)
		return nil, err
	} else if s.Data["password"] == nil {
		err := errors.Errorf("credentials: could not load credentials from %s", cfg.ID)
		span.RecordError(err)
		return nil, err
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
