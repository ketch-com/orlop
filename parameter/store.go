package parameter

import (
	"context"
	"go.ketch.com/lib/orlop/errors"
)

var ErrNotFound = errors.New("not found")

// Store provide an interface to interact with a parameter store
type Store interface {
	List(ctx context.Context, p string) ([]string, error)
	Read(ctx context.Context, p string) (map[string]interface{}, error)
	Write(ctx context.Context, p string, data map[string]interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context, p string) error
}

// ObjectStore is an extension to read/write objects to parameter store
type ObjectStore interface {
	Store

	ReadObject(ctx context.Context, p string, out interface{}) error
	WriteObject(ctx context.Context, p string, in interface{}) error
}

// StoreFromObjectStore returns a Store given an ObjectStore
func StoreFromObjectStore(store ObjectStore) Store {
	return store
}
