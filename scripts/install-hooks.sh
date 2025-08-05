#!/bin/sh

# Script to install Git hooks

echo "Installing Git pre-commit hook..."

# Create scripts directory if it doesn't exist
mkdir -p scripts

# Copy pre-commit hook
cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "Git pre-commit hook installed successfully!"
echo "The hook will run golangci-lint before each commit."
