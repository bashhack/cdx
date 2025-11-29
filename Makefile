# ============================================================================= #
# HELPERS
# ============================================================================= #

## help: Print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## search: Fuzzy search and run commands (requires fzf)
.PHONY: search
search:
	@if ! command -v fzf >/dev/null 2>&1; then \
		echo "❌ Error: 'fzf' is not installed."; \
		echo "Install with: brew install fzf"; \
		exit 1; \
	fi
	@target=$$(sed -n 's/^##//p' ${MAKEFILE_LIST} | fzf --height=50% --reverse --header="Select a command to run" | awk '{print $$1}'); \
	if [ -n "$$target" ]; then \
		echo "Running: make $$target"; \
		$(MAKE) $$target; \
	fi

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ============================================================================= #
# DEVELOPMENT
# ============================================================================= #

## run: Run the application
.PHONY: run
run:
	go run ./cmd/cdx

## run/tests: Run test suite
.PHONY: run/tests
run/tests:
	go test -v ./...

## dev/setup/hooks: Install git hooks for pre-commit and commit-msg checks
.PHONY: dev/setup/hooks
dev/setup/hooks:
	@echo 'Installing git hooks...'
	@if [ ! -f .githooks/pre-commit ]; then \
		echo '❌ Error: .githooks/pre-commit not found'; \
		exit 1; \
	fi
	@if [ ! -f .githooks/commit-msg ]; then \
		echo '❌ Error: .githooks/commit-msg not found'; \
		exit 1; \
	fi
	@mkdir -p .git/hooks
	@cp .githooks/pre-commit .git/hooks/pre-commit
	@cp .githooks/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/pre-commit
	@chmod +x .git/hooks/commit-msg
	@echo '✅ Git hooks installed (pre-commit + commit-msg)'
	@echo ''
	@echo 'Commit messages must follow Conventional Commits format:'
	@echo '  <type>[optional scope]: <description>'
	@echo ''
	@echo 'Allowed types: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test, wip'

# ============================================================================= #
# QUALITY CONTROL
# ============================================================================= #

# Internal helper - not in help menu
.PHONY: check_staticcheck
check_staticcheck:
	@if ! command -v staticcheck >/dev/null 2>&1; then \
		echo "Installing staticcheck..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi

## format: Format all Go code with goimports
.PHONY: format
format:
	@echo 'Formatting Go code...'
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@goimports -w -local github.com/bashhack/cdx $$(find . -name '*.go' -not -path "./vendor/*")
	@echo '✅ Code formatted'

## format/check: Check if code is properly formatted (non-destructive)
.PHONY: format/check
format/check:
	@echo 'Checking Go code formatting...'
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@if [ -n "$$(goimports -l -local github.com/bashhack/cdx $$(find . -name '*.go' -not -path './vendor/*'))" ]; then \
		echo "❌ The following files need formatting:"; \
		goimports -l -local github.com/bashhack/cdx $$(find . -name '*.go' -not -path './vendor/*'); \
		echo "Run 'make format' to fix"; \
		exit 1; \
	fi
	@echo '✅ All files properly formatted'

## lint: Run linters without tests
.PHONY: lint
lint:
	@echo 'Formatting code...'
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@goimports -w -local github.com/bashhack/cdx $$(find . -name '*.go' -not -path "./vendor/*")
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running staticcheck...'
	@$(MAKE) check_staticcheck
	staticcheck ./...
	@echo '✅ Linting complete'

## lint/golangci: Run golangci-lint (comprehensive linting tool)
.PHONY: lint/golangci
lint/golangci:
	@echo 'Running golangci-lint...'
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint v2.6.2..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.6.2; \
	fi
	@golangci-lint run ./...
	@echo '✅ golangci-lint complete'

## audit: Tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		go install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@goimports -w -local github.com/bashhack/cdx $$(find . -name '*.go' -not -path "./vendor/*")
	@echo 'Vetting code...'
	go vet ./...
	@echo 'Running staticcheck...'
	@$(MAKE) check_staticcheck
	staticcheck ./...
	@echo 'Running tests...'
	go test -short -vet=off ./...
	@echo '✅ Audit complete'

## vendor: Tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

## coverage: Run test suite with coverage
.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# ============================================================================= #
# BUILD
# ============================================================================= #

## build: Build the application
.PHONY: build
build:
	@echo 'Building cdx...'
	go build -o=./cdx ./cmd/cdx

## build/optimize: Build optimized application (sans DWARF + symbol table)
.PHONY: build/optimize
build/optimize:
	@echo 'Building optimized cdx...'
	go build -ldflags='-s -w' -o=./cdx ./cmd/cdx

## install: Install cdx locally
.PHONY: install
install:
	@echo 'Installing cdx...'
	go install ./cmd/cdx
