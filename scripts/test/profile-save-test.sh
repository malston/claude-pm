#!/bin/bash
# Tests profile save functionality

set -e
echo "=== Profile Save Tests ==="

# Test 1: Save to new profile
echo "[1] Testing: save to new profile name"
claudeup profile save test-save-new
claudeup profile show test-save-new
echo "✓ Pass"

# Test 2: Save to existing profile (will prompt - tester confirms)
echo "[2] Testing: save to existing profile (answer 'y' when prompted)"
claudeup profile save test-save-new
echo "✓ Pass"

# Test 3: Save with -y skips prompt
echo "[3] Testing: save with -y flag"
claudeup profile save test-save-new -y
echo "✓ Pass"

# Test 4: Save with no name uses active profile
echo "[4] Testing: save with no name (uses active profile)"
claudeup profile use test-save-new -y
claudeup profile save
echo "✓ Pass"

# Cleanup
rm -f ~/.claudeup/profiles/test-save-new.json
echo "=== All profile save tests passed ==="
