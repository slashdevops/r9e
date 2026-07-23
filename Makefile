.DELETE_ON_ERROR: clean

EXECUTABLES = go
K := $(foreach exec,$(EXECUTABLES),\
  $(if $(shell command -v $(exec)),some string,$(error "No $(exec) in PATH")))

# Ensure shell errors are propagated.
.SHELLFLAGS := -e -c

PROJECT_NAME      ?= $(shell grep '^module' go.mod | cut -d '/' -f 3)
PROJECT_NAMESPACE ?= $(shell grep '^module' go.mod | cut -d '/' -f 2)
PROJECT_DEPENDENCIES := $(shell go list -m -f '{{if not (or .Indirect .Main)}}{{.Path}}{{end}}' all)

BUILD_DIR ?= ./build

PROJECT_COVERAGE_FILE ?= $(BUILD_DIR)/coverage.txt
PROJECT_COVERAGE_MODE ?= atomic

######## Functions ########
# exec_cmd runs a command with friendly output.
#   MAKE_DEBUG=true          print the command instead of running it.
#   MAKE_STOP_ON_ERRORS=true abort the whole run when a command fails (useful in CI).
MAKE_STOP_ON_ERRORS ?= false
MAKE_DEBUG          ?= false

define exec_cmd
$(if $(filter $(MAKE_DEBUG),true),\
	${1} \
, \
	$(if $(filter $(MAKE_STOP_ON_ERRORS),true),\
		$(if $(findstring >, $1),\
			@${1} 2>/dev/null && printf "  🤞 ${1} ✅\n" || (printf "  ${1} ❌\n"; exit 1) \
		, \
			@${1} > /dev/null && printf "  🤞 ${1} ✅\n" || (printf "  ${1} ❌\n"; exit 1) \
		) \
	, \
		$(if $(findstring >, $1),\
			@${1} 2>/dev/null; _exit_code=$$?; if [ $$_exit_code -eq 0 ]; then printf "  🤞 ${1} ✅\n"; else printf "  ${1} ❌\n"; fi; exit $$_exit_code \
		, \
			@${1} > /dev/null 2>&1; _exit_code=$$?; if [ $$_exit_code -eq 0 ]; then printf '  🤞 ${1} ✅\n'; else printf '  ${1} ❌\n'; fi; exit $$_exit_code \
		) \
	) \
)

endef # don't remove the white line before endef

###############################################################################
######## Targets ##############################################################
##@ Default command
.PHONY: all
all: clean check ## Clean and run the full quality gate (default target).

###############################################################################
##@ Golang commands
.PHONY: go-fmt
go-fmt: ## Format go code.
	@printf "👉 Formatting go code...\n"
	$(call exec_cmd, go fmt ./... )

.PHONY: go-vet
go-vet: ## Vet go code.
	@printf "👉 Vet go code...\n"
	$(call exec_cmd, go vet ./... )

.PHONY: go-fix
go-fix: go-fmt go-vet ## Apply modernizations, then fmt and vet.
	@printf "👉 Fixing go code...\n"
	$(call exec_cmd, go fix ./... )

.PHONY: go-betteralign
go-betteralign: install-betteralign ## Align struct fields for optimal memory layout.
	@printf "👉 Aligning struct fields with betteralign...\n"
	$(call exec_cmd, betteralign -apply ./... )

.PHONY: go-mod-tidy
go-mod-tidy: ## Clean go.mod and go.sum.
	@printf "👉 Cleaning go.mod and go.sum...\n"
	$(call exec_cmd, go mod tidy)

.PHONY: go-mod-update
go-mod-update: go-mod-tidy ## Update all direct dependencies.
	@printf "👉 Updating dependencies...\n"
	$(foreach DEP, $(PROJECT_DEPENDENCIES), \
		$(call exec_cmd, go get -u $(DEP)) \
	)
	$(call exec_cmd, go mod tidy)

.PHONY: go-mod-verify
go-mod-verify: ## Verify go.mod and go.sum.
	@printf "👉 Verifying modules...\n"
	$(call exec_cmd, go mod verify)

###############################################################################
##@ Test commands
$(PROJECT_COVERAGE_FILE):
	@printf "👉 Creating coverage file...\n"
	$(call exec_cmd, mkdir -p $(BUILD_DIR) )
	$(call exec_cmd, touch $(PROJECT_COVERAGE_FILE) )

.PHONY: test
test: $(PROJECT_COVERAGE_FILE) ## Run tests with the race detector and coverage.
	@printf "👉 Running tests...\n"
	$(call exec_cmd, go test \
		-race \
		-coverprofile=$(PROJECT_COVERAGE_FILE) \
		-covermode=$(PROJECT_COVERAGE_MODE) \
		./... \
	)

.PHONY: test-coverage
test-coverage: test ## Open the HTML coverage report in the browser.
	@printf "👉 Opening coverage report...\n"
	$(call exec_cmd, go tool cover -html=$(PROJECT_COVERAGE_FILE))

.PHONY: cover-func
cover-func: test ## Print per-function and total coverage.
	@printf "👉 Coverage summary...\n"
	$(call exec_cmd, go tool cover -func=$(PROJECT_COVERAGE_FILE))

.PHONY: bench
bench: ## Run all benchmarks with allocation stats.
	@printf "👉 Running benchmarks...\n"
	$(call exec_cmd, go test -run '^$$' -bench . -benchmem ./... )

###############################################################################
##@ Build & Check commands
.PHONY: build
build: ## Verify the package builds.
	@printf "👉 Building...\n"
	$(call exec_cmd, go build ./... )

.PHONY: lint
lint: install-golangci-lint ## Lint go code with golangci-lint.
	@printf "👉 Linting...\n"
	$(call exec_cmd, golangci-lint run ./... )

.PHONY: vulncheck
vulncheck: install-govulncheck ## Check for known vulnerabilities.
	@printf "👉 Checking vulnerabilities...\n"
	$(call exec_cmd, govulncheck ./... )

.PHONY: check
check: go-fix test build ## Run the local quality gate (matches CI): fix, test, build.

###############################################################################
##@ Tools commands
.PHONY: install-golangci-lint
install-golangci-lint: ## Install golangci-lint (https://golangci-lint.run/).
	@printf "👉 Installing golangci-lint...\n"
	$(call exec_cmd, go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2 )

.PHONY: install-betteralign
install-betteralign: ## Install betteralign (https://github.com/dkorunic/betteralign).
	@printf "👉 Installing betteralign...\n"
	$(call exec_cmd, go install github.com/dkorunic/betteralign/cmd/betteralign@latest )

.PHONY: install-govulncheck
install-govulncheck: ## Install govulncheck (https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck).
	@printf "👉 Installing govulncheck...\n"
	$(call exec_cmd, go install golang.org/x/vuln/cmd/govulncheck@latest )

###############################################################################
##@ Support commands
.PHONY: clean
clean: ## Remove build and coverage artifacts.
	@printf "👉 Cleaning environment...\n"
	$(call exec_cmd, rm -rf $(BUILD_DIR) ./*.out )

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##";                                             \
		printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ \
		{ printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/            \
		{ printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } '                  \
		$(MAKEFILE_LIST)
