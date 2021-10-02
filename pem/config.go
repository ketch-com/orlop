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

package pem

const (
	// CertificateWhich is the key of the TLS certificate in secrets engine
	CertificateWhich = "certificate"
	// KeyWhich is the key of the TLS private key in secrets engine
	KeyWhich = "key"
	// RootCAWhich is the key of the TLS root CA in secrets engine
	RootCAWhich = "issuing_ca"
	// PublicKeyWhich is the key of the public key in secrets engine
	PublicKeyWhich = "public_key"
	// PrivateKeyWhich is the key of the private key in secrets engine
	PrivateKeyWhich = "private_key"
)

// Config provides pem-related configurations
type Config interface {
	GetID() string
	GetFile() string
	GetEnabled() bool
	GetWhich() string
}

// CertificateConfig provides pem-related configurations
type CertificateConfig struct {
	ID   string
	File string
}

func CertificateID(id string) Config {
	return CertificateConfig{
		ID: id,
	}
}

func CertificateFile(file string) Config {
	return CertificateConfig{
		File: file,
	}
}

func (c CertificateConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.File) > 0
}

func (c CertificateConfig) GetID() string {
	return c.ID
}

func (c CertificateConfig) GetFile() string {
	return c.File
}

func (c CertificateConfig) GetWhich() string {
	return CertificateWhich
}

// KeyConfig provides pem-related configurations
type KeyConfig struct {
	ID   string
	File string
}

func KeyID(id string) Config {
	return KeyConfig{
		ID: id,
	}
}

func KeyFile(file string) Config {
	return KeyConfig{
		File: file,
	}
}

func (c KeyConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.File) > 0
}

func (c KeyConfig) GetID() string {
	return c.ID
}

func (c KeyConfig) GetFile() string {
	return c.File
}

func (c KeyConfig) GetWhich() string {
	return KeyWhich
}

// RootCAConfig provides pem-related configurations
type RootCAConfig struct {
	ID   string
	File string
}

func RootCAID(id string) Config {
	return RootCAConfig{
		ID: id,
	}
}

func RootCAFile(file string) Config {
	return RootCAConfig{
		File: file,
	}
}

func (c RootCAConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.File) > 0
}

func (c RootCAConfig) GetID() string {
	return c.ID
}

func (c RootCAConfig) GetFile() string {
	return c.File
}

func (c RootCAConfig) GetWhich() string {
	return RootCAWhich
}

// PublicKeyConfig provides pem-related configurations
type PublicKeyConfig struct {
	ID   string
	File string
}

func PublicKeyID(id string) Config {
	return PublicKeyConfig{
		ID: id,
	}
}

func PublicKeyFile(file string) Config {
	return PublicKeyConfig{
		File: file,
	}
}

func (c PublicKeyConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.File) > 0
}

func (c PublicKeyConfig) GetID() string {
	return c.ID
}

func (c PublicKeyConfig) GetFile() string {
	return c.File
}

func (c PublicKeyConfig) GetWhich() string {
	return PublicKeyWhich
}

// PrivateKeyConfig provides pem-related configurations
type PrivateKeyConfig struct {
	ID   string
	File string
}

func PrivateKeyID(id string) Config {
	return PrivateKeyConfig{
		ID: id,
	}
}

func PrivateKeyFile(file string) Config {
	return PrivateKeyConfig{
		File: file,
	}
}

func (c PrivateKeyConfig) GetEnabled() bool {
	return len(c.ID) > 0 || len(c.File) > 0
}

func (c PrivateKeyConfig) GetID() string {
	return c.ID
}

func (c PrivateKeyConfig) GetFile() string {
	return c.File
}

func (c PrivateKeyConfig) GetWhich() string {
	return PrivateKeyWhich
}
