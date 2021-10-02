package health

import "context"

// Check provides the capability to check the health
type Check struct {
	Name    string
	Checker func(ctx context.Context) (interface{}, error)
}

