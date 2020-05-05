package orlop

import "os"

// EnvironmentKey is the environment variable we look for to set the environment
const EnvironmentKey = "SWITCHBIT_ENVIRONMENT"

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
