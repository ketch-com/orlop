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

import "os"

// EnvironmentKey is the environment variable we look for to set the environment
const EnvironmentKey = "KETCH_ENVIRONMENT"

// Environment is a defined environment
type Environment string

// IsLocal returns true if the environment is not defined (aka local)
func (e Environment) IsLocal() bool {
	return e == ""
}

// IsProduction returns true if the environment is the production environment.
func (e Environment) IsProduction() bool {
	return e == "prod" || e == "production"
}

// IsTest returns true if the environment is the test environment
func (e Environment) IsTest() bool {
	return e == "test"
}

// String returns a string version of the environment.
func (e Environment) String() string {
	return string(e)
}

// Env returns the environment from the environment variables
func Env() Environment {
	return Environment(os.Getenv(EnvironmentKey))
}
