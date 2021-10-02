package parameter

import (
	"context"
)

type noopStore struct{}

func NewNoopStore() ObjectStore {
	return &noopStore{}
}

func (c noopStore) List(ctx context.Context, p string) ([]string, error) {
	return nil, nil
}

func (c noopStore) Read(ctx context.Context, p string) (map[string]interface{}, error) {
	return nil, ErrNotFound
}

func (c noopStore) Write(ctx context.Context, p string, data map[string]interface{}) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

func (c noopStore) Delete(ctx context.Context, p string) error {
	return nil
}

func (c noopStore) ReadObject(ctx context.Context, p string, out interface{}) error {
	return nil
}

func (c noopStore) WriteObject(ctx context.Context, p string, in interface{}) error {
	return nil
}

func (c noopStore) GetEnabled() bool {
	return false
}
