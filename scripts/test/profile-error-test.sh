#!/bin/bash
# Tests error conditions

echo "=== Profile Error Tests ==="

# Test 1: create -y with no active profile
echo "[1] Testing: create -y with no active profile errors"
# Clear active profile first (may fail if no config exists - that's ok)
claudeup config set preferences.activeProfile "" 2>/dev/null || true
if claudeup profile create test-err -y 2>&1 | grep -qi "no active profile"; then
  echo "✓ Pass"
else
  echo "✗ Fail - should have errored about no active profile"
  exit 1
fi

# Test 2: create --from nonexistent profile
echo "[2] Testing: create --from nonexistent profile errors"
if claudeup profile create test-err --from nonexistent-xyz 2>&1 | grep -q "not found"; then
  echo "✓ Pass"
else
  echo "✗ Fail - should have errored about profile not found"
  exit 1
fi

echo "=== All error tests passed ==="
