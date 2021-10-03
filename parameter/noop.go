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
