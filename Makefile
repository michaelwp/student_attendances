migration-up:
	@echo "Running database migrations up..."
	@go run db/migration.go up

migration-down:
	@echo "Running database migrations down..."
	@go run db/migration.go down

install-hooks: ## Install Git hooks
	@echo "Installing Git pre-commit hook..."
	@mkdir -p scripts
	@cp scripts/pre-commit.sh .git/hooks/pre-commit 2>/dev/null || echo "#!/bin/sh\n\n# Git pre-commit hook to run golangci-lint\n\n# Colors for output\nRED='\\033[0;31m'\nGREEN='\\033[0;32m'\nYELLOW='\\033[1;33m'\nNC='\\033[0m' # No Color\n\necho \"\$${YELLOW}Running golangci-lint...\$${NC}\"\n\n# Check if golangci-lint is installed\nif ! command -v golangci-lint &> /dev/null; then\n    echo \"\$${RED}Error: golangci-lint is not installed.\$${NC}\"\n    echo \"Please install it by running:\"\n    echo \"  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest\"\n    echo \"Or visit: https://golangci-lint.run/usage/install/\"\n    exit 1\nfi\n\n# Run golangci-lint on staged files\nif ! golangci-lint run --config .golangci.yml; then\n    echo \"\$${RED}golangci-lint found issues. Please fix them before committing.\$${NC}\"\n    echo \"You can run the following command to see the issues:\"\n    echo \"  golangci-lint run\"\n    exit 1\nfi\n\necho \"\$${GREEN}golangci-lint passed! ✓\$${NC}\"\nexit 0" > .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Git pre-commit hook installed successfully!"

uninstall-hooks: ## Remove Git hooks
	@echo "Removing Git pre-commit hook..."
	@rm -f .git/hooks/pre-commit
	@echo "Git pre-commit hook removed!"

lint: ## Run golangci-lint
	golangci-lint run --config .golangci.yml

lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --config .golangci.yml --fix

swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/student_attendance/main.go --output docs --parseDependency --parseInternal
	@echo "Swagger documentation generated successfully!"

build: ## Build the application
	@echo "Building the application..."
	@go build -o bin/student_attendance cmd/student_attendance/main.go

run: ## Run the application
	@echo "Running the application..."
	@./bin/student_attendance

.PHONY: migration-up migration-down install-hooks ui uninstall-hooks lint lint-fix swagger build run