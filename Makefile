COMPLEXITY_THRESHOLD := 15
COVERAGE_THRESHOLD := 90
COVERAGE_PACKAGES := ./adapter/...

.PHONY: all build test test-race test-tui test-tui-update lint lint-complexity lint-coverage check fmt vet coverage-html vhs-record vhs-validate

all: check build

build:
	go build ./...

test:
	go test ./...

test-race:
	go test -race ./...

test-tui:
	go test ./... -tags=teatest

test-tui-update:
	go test ./... -tags=teatest -args -update

lint-complexity:
	@echo "Checking complexity (threshold: $(COMPLEXITY_THRESHOLD))..."
	@result=$$(gocyclo -over $(COMPLEXITY_THRESHOLD) -ignore "example/.*" .); \
	if [ -n "$$result" ]; then \
		echo "$$result"; \
		echo "ERROR: Functions exceed complexity threshold"; \
		exit 1; \
	fi
	@echo "OK: All functions within complexity threshold"

lint-coverage:
	@echo "Checking coverage (threshold: $(COVERAGE_THRESHOLD)%)..."
	@go test -cover $(COVERAGE_PACKAGES) -coverprofile=coverage.out > /dev/null 2>&1
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | tr -d '%'); \
	if [ $$(echo "$$coverage < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "ERROR: Coverage $$coverage% < $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	fi
	@echo "OK: Coverage meets threshold"
	@rm -f coverage.out

lint: lint-complexity lint-coverage

check: fmt vet lint test

fmt:
	go fmt ./...

vet:
	go vet ./...

coverage-html:
	go test -cover $(COVERAGE_PACKAGES) -coverprofile=coverage.out
	go tool cover -html=coverage.out

# VHS visual regression testing
vhs-record:
	@echo "Recording VHS screenshots..."
	@mkdir -p example/chatinput/screenshots
	vhs example/chatinput/chatinput.tape
	@echo "Screenshots saved to example/chatinput/screenshots/"

vhs-validate:
	@echo "Validating VHS tape files..."
	vhs validate example/chatinput/chatinput.tape
	@echo "OK: Tape file syntax is valid"
