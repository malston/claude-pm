#!/bin/bash
# Tests profile create (copy/fork) functionality

set -e
echo "=== Profile Create Tests ==="

# Test 1: Create with --from
echo "[1] Testing: create with --from flag"
claudeup profile create test-fork-1 --from default
claudeup profile show test-fork-1
echo "✓ Pass"

# Test 2: Create with -y uses active profile
echo "[2] Testing: create with -y uses active profile"
claudeup profile use default -y
claudeup profile create test-fork-2 -y
claudeup profile show test-fork-2
echo "✓ Pass"

# Test 3: Create when name exists should error
echo "[3] Testing: create existing name errors"
if claudeup profile create test-fork-1 --from default 2>&1 | grep -q "already exists"; then
  echo "✓ Pass"
else
  echo "✗ Fail - should have errored"
  exit 1
fi

# Test 4: Interactive picker (tester selects manually)
echo "[4] Testing: interactive picker (select a profile when prompted)"
claudeup profile create test-fork-3
claudeup profile show test-fork-3
echo "✓ Pass"

# Cleanup
rm -f ~/.claudeup/profiles/test-fork-*.json
echo "=== All profile create tests passed ==="
