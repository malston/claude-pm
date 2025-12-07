// ABOUTME: Tests for environment variable resolver
// ABOUTME: Validates reading secrets from environment
package secrets

import (
	"os"
	"testing"
)

func TestEnvResolverName(t *testing.T) {
	r := NewEnvResolver()
	if r.Name() != "env" {
		t.Errorf("Expected name 'env', got %q", r.Name())
	}
}

func TestEnvResolverAlwaysAvailable(t *testing.T) {
	r := NewEnvResolver()
	if !r.Available() {
		t.Error("EnvResolver should always be available")
	}
}

func TestEnvResolverSuccess(t *testing.T) {
	r := NewEnvResolver()

	// Set a test env var
	testKey := "CLAUDE_PM_TEST_SECRET"
	testValue := "test-value-12345"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)

	value, err := r.Resolve(testKey)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if value != testValue {
		t.Errorf("Expected %q, got %q", testValue, value)
	}
}

func TestEnvResolverMissing(t *testing.T) {
	r := NewEnvResolver()

	_, err := r.Resolve("DEFINITELY_NOT_SET_CLAUDE_PM_TEST")
	if err == nil {
		t.Error("Expected error for missing env var")
	}
}
