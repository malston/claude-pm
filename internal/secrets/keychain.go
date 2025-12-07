// ABOUTME: macOS Keychain secret resolver
// ABOUTME: Uses 'security find-generic-password' to fetch secrets
package secrets

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

// KeychainResolver resolves secrets from macOS Keychain
type KeychainResolver struct {
	available *bool
}

// NewKeychainResolver creates a new Keychain resolver
func NewKeychainResolver() *KeychainResolver {
	return &KeychainResolver{}
}

// Name returns the resolver identifier
func (k *KeychainResolver) Name() string {
	return "keychain"
}

// Available returns true if running on macOS
func (k *KeychainResolver) Available() bool {
	if k.available != nil {
		return *k.available
	}

	available := runtime.GOOS == "darwin"
	k.available = &available
	return available
}

// Resolve fetches a secret from macOS Keychain
// ref should be in the format: service:account or just service
func (k *KeychainResolver) Resolve(ref string) (string, error) {
	parts := strings.SplitN(ref, ":", 2)
	service := parts[0]
	account := ""
	if len(parts) > 1 {
		account = parts[1]
	}

	args := []string{"find-generic-password", "-s", service, "-w"}
	if account != "" {
		args = append(args, "-a", account)
	}

	cmd := exec.Command("security", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
