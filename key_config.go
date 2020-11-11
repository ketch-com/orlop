// Copyright (c) 2020 Ketch, Inc.
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

// HasKeyConfig denotes that the object provides Key configuration
type HasKeyConfig interface {
	GetID() string
	GetSecret() []byte
	GetFile() string
}

// KeyConfig provides key-related configurations
type KeyConfig struct {
	ID     string
	Secret []byte `config:"secret,encoding=base64"`
	File   string
}

// GetID returns the Vault path to retrieve the key from
func (c KeyConfig) GetID() string {
	return c.ID
}

// GetSecret returns the secret (most likely from an environment variable)
func (c KeyConfig) GetSecret() []byte {
	return c.Secret
}

// GetFile returns the path to a file containing the key
func (c KeyConfig) GetFile() string {
	return c.File
}

// GetEnabled returns true if the key is enabled
func (c KeyConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.Secret) > 0 || len(c.File) > 0
}
