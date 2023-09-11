package log

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

// Publisher is an implementation of a Publisher using Log
type Publisher struct{}

// NewPublisher returns a new Publisher
func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) PublishEvent(ctx context.Context, subject string, event any) error {
	if event == nil {
		return nil
	}

	fields := make(logrus.Fields)
	if err := mapstructure.Decode(event, &fields); err != nil {
		return err
	}

	WithContext(ctx).WithFields(logrus.Fields{
		"subject": subject,
		"data":    fields,
	}).Info("published")

	return nil
}
