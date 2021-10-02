package secret

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.ketch.com/lib/orlop/config"
	"go.ketch.com/lib/orlop/errors"
	"go.ketch.com/lib/orlop/logging"
	"go.ketch.com/lib/orlop/parameter"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"io/ioutil"
)

type Provider interface {
	Load(ctx context.Context, cfg Config) (string, error)
}

type Params struct {
	fx.In

	Secrets parameter.Store
	Logger  logging.Logger
	Tracer  trace.Tracer
}

func New(params Params) Provider {
	return &providerImpl{
		params: params,
	}
}

type providerImpl struct {
	params Params
}

// Load the secret material from the secrets engine based on the config
func (p *providerImpl) Load(ctx context.Context, cfg Config) (string, error) {
	// If the secret is not enabled, return an empty string
	if !cfg.GetEnabled() {
		return "", nil
	}

	ctx, span := p.params.Tracer.Start(ctx, "LoadSecret")
	defer span.End()

	fields := logrus.Fields{
		"which": cfg.Which,
	}

	method := "none"

	if len(cfg.ID) > 0 {
		if config.IsEnabled(p.params.Secrets) {
			method = "id"
		}
		fields["secret.id"] = cfg.ID
		span.SetAttributes(attribute.String("secret.id", cfg.ID))
	}

	if len(cfg.File) > 0 {
		method = "file"
		fields["secret.file"] = cfg.File
		span.SetAttributes(attribute.String("secret.file", cfg.File))
	}

	if len(cfg.Secret) > 0 {
		method = "value"
		fields["secret.value"] = "*********"
		span.SetAttributes(attribute.String("secret.value", "*********"))
	}

	fields["method"] = method
	span.SetAttributes(attribute.String("secret.method", method))
	l := p.params.Logger.WithFields(fields)

	switch method {
	case "value":
		l.Trace("secret found")
		return string(cfg.Secret), nil

	case "file":
		secretBytes, err := ioutil.ReadFile(cfg.File)
		if err != nil {
			err = errors.Wrap(err, "pem: not found")
			span.RecordError(err)
			return "", err
		}

		return string(secretBytes), nil

	case "id":
		s, err := p.params.Secrets.Read(ctx, cfg.ID)
		if err != nil {
			err = errors.Wrap(err, "secret: not found")
			span.RecordError(err)
			return "", err
		}

		if s == nil || s[cfg.Which] == nil {
			err = errors.New("secret: not found")
			span.RecordError(err)
			return "", err
		}

		l.Trace("secret found")
		return s[cfg.Which].(string), nil
	}

	err := errors.Errorf("secret: no secret configured for %s", cfg.Which)
	span.RecordError(err)
	return "", err
}
