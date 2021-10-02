// Copyright (c) 2021 Ketch Kloud, Inc.
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
