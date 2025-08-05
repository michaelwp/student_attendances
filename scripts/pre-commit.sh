
#!/bin/sh

# Git pre-commit hook to run golangci-lint

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "${YELLOW}Running golangci-lint...${NC}"

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo "${RED}Error: golangci-lint is not installed.${NC}"
    echo "Please install it by running:"
    echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    echo "Or visit: https://golangci-lint.run/usage/install/"
    exit 1
fi

# Run golangci-lint on staged files
if ! golangci-lint run --config .golangci.yml; then
    echo "${RED}golangci-lint found issues. Please fix them before committing.${NC}"
    echo "You can run the following command to see the issues:"
    echo "  golangci-lint run"
    exit 1
fi

echo "${GREEN}golangci-lint passed! ✓${NC}"
exit 0
