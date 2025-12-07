// ABOUTME: Environment variable secret resolver
// ABOUTME: Reads secrets from environment variables
package secrets

import (
	"errors"
	"os"
)

// EnvResolver resolves secrets from environment variables
type EnvResolver struct{}

// NewEnvResolver creates a new environment variable resolver
func NewEnvResolver() *EnvResolver {
	return &EnvResolver{}
}

// Name returns the resolver identifier
func (e *EnvResolver) Name() string {
	return "env"
}

// Available always returns true - env vars are always accessible
func (e *EnvResolver) Available() bool {
	return true
}

// Resolve gets the value of the specified environment variable
func (e *EnvResolver) Resolve(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("environment variable not set: " + key)
	}
	return value, nil
}
