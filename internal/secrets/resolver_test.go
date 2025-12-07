// ABOUTME: Tests for secret resolution chain
// ABOUTME: Validates multi-backend resolution with priority ordering
package secrets

import (
	"errors"
	"testing"
)

// mockResolver is a test resolver
type mockResolver struct {
	name      string
	available bool
	value     string
	err       error
}

func (m *mockResolver) Name() string       { return m.name }
func (m *mockResolver) Available() bool    { return m.available }
func (m *mockResolver) Resolve(_ string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.value, nil
}

func TestChainResolvesFirstAvailable(t *testing.T) {
	chain := NewChain(
		&mockResolver{name: "first", available: false, value: "wrong"},
		&mockResolver{name: "second", available: true, value: "correct"},
		&mockResolver{name: "third", available: true, value: "also-wrong"},
	)

	value, source, err := chain.Resolve("any-ref")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if value != "correct" {
		t.Errorf("Expected 'correct', got %q", value)
	}

	if source != "second" {
		t.Errorf("Expected source 'second', got %q", source)
	}
}

func TestChainSkipsErrors(t *testing.T) {
	chain := NewChain(
		&mockResolver{name: "failing", available: true, err: errors.New("failed")},
		&mockResolver{name: "working", available: true, value: "success"},
	)

	value, source, err := chain.Resolve("any-ref")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if value != "success" {
		t.Errorf("Expected 'success', got %q", value)
	}

	if source != "working" {
		t.Errorf("Expected source 'working', got %q", source)
	}
}

func TestChainErrorsWhenAllFail(t *testing.T) {
	chain := NewChain(
		&mockResolver{name: "first", available: false},
		&mockResolver{name: "second", available: true, err: errors.New("failed")},
	)

	_, _, err := chain.Resolve("any-ref")
	if err == nil {
		t.Error("Expected error when all resolvers fail")
	}
}

func TestEmptyChainErrors(t *testing.T) {
	chain := NewChain()

	_, _, err := chain.Resolve("any-ref")
	if err == nil {
		t.Error("Expected error for empty chain")
	}
}
