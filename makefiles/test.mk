# Testing Targets
# This file contains targets for testing, coverage, and quality assurance

.PHONY: test
test: ## Run tests
	go test $(TEST_FLAGS) ./...

.PHONY: test-short
test-short: ## Run tests in short mode
	go test $(TEST_FLAGS) -short ./...

.PHONY: test-race
test-race: ## Run tests with race detection
	go test $(TEST_FLAGS) -race ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	go test $(TEST_FLAGS) -coverprofile=$(COVERAGE_OUTPUT) ./...
	go tool cover -html=$(COVERAGE_OUTPUT) -o $(COVERAGE_HTML)

.PHONY: test-cover-func
test-cover-func: ## Run tests with function coverage
	go test $(TEST_FLAGS) -coverprofile=$(COVERAGE_OUTPUT) ./...
	go tool cover -func=$(COVERAGE_OUTPUT)

.PHONY: test-bench
test-bench: ## Run benchmark tests
	go test $(TEST_FLAGS) -bench=. -benchmem ./...

.PHONY: test-bench-cpu
test-bench-cpu: ## Run CPU profile benchmark
	go test $(TEST_FLAGS) -bench=. -cpuprofile=cpu.prof ./...

.PHONY: test-bench-mem
test-bench-mem: ## Run memory profile benchmark
	go test $(TEST_FLAGS) -bench=. -memprofile=mem.prof ./...

.PHONY: test-integration
test-integration: ## Run integration tests
	go test $(TEST_FLAGS) -tags=integration ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	go test $(TEST_FLAGS) -short ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	go test -v ./...

.PHONY: test-coverage-view
test-coverage-view: ## View coverage report in browser
	@if [ -f $(COVERAGE_HTML) ]; then \
		open $(COVERAGE_HTML); \
	else \
		echo "Coverage report not found. Run 'make test-cover' first."; \
	fi

.PHONY: test-clean
test-clean: ## Clean test artifacts
	rm -f $(COVERAGE_OUTPUT) $(COVERAGE_HTML)
	rm -f *.prof
	rm -f test.log
