// ABOUTME: 1Password CLI secret resolver
// ABOUTME: Uses 'op read' to fetch secrets from 1Password vaults
package secrets

import (
	"bytes"
	"os/exec"
	"strings"
)

// OnePasswordResolver resolves secrets using 1Password CLI
type OnePasswordResolver struct {
	available *bool
}

// NewOnePasswordResolver creates a new 1Password resolver
func NewOnePasswordResolver() *OnePasswordResolver {
	return &OnePasswordResolver{}
}

// Name returns the resolver identifier
func (o *OnePasswordResolver) Name() string {
	return "1password"
}

// Available returns true if the 'op' CLI is installed
func (o *OnePasswordResolver) Available() bool {
	if o.available != nil {
		return *o.available
	}

	_, err := exec.LookPath("op")
	available := err == nil
	o.available = &available
	return available
}

// Resolve fetches a secret from 1Password using 'op read'
// ref should be in the format: op://vault/item/field
func (o *OnePasswordResolver) Resolve(ref string) (string, error) {
	cmd := exec.Command("op", "read", ref)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
