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

package env

import (
	"github.com/joho/godotenv"
	"os"
)

// EnvironmentKey is the environment variable we look for to set the environment (prefixed with service name)
const EnvironmentKey = "ENVIRONMENT"

// Environment is a defined environment
type Environment string

// Env returns the environment from the environment variables
func Env(environ Environ) Environment {
	return Environment(environ.Getenv(EnvironmentKey))
}

// IsLocal returns true if the environment is not defined (aka local)
func (e Environment) IsLocal() bool {
	return e == "" || e == "local"
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

// Load loads the environment variables from the specified files and from the standard locations
func (e Environment) Load(files ...string) {
	var envFiles []string
	envFiles = append(envFiles, files...)
	if e.IsLocal() {
		envFiles = append(envFiles, ".env.local")
	} else {
		envFiles = append(envFiles, ".env."+e.String())
	}
	envFiles = append(envFiles, ".env")

	for _, file := range envFiles {
		if _, err := os.Stat(file); err == nil {
			_ = godotenv.Load(file)
		}
	}
}
