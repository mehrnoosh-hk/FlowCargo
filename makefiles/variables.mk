# FlowCargo Variables
# This file contains all shared variables used across different Makefiles

# Application Configuration
BINARY_NAME=flowcargo
MAIN_PATH=./cmd/api
MIGRATION_DIR=./migrations

# Load environment variables from .env.dev if it exists
ifneq (,$(wildcard .env.dev))
    include .env.dev
    export
endif

# Default migration DSN if not set in environment
MIGRATION_DSN ?= postgres://user:password@localhost/dbname?sslmode=disable

# Go Configuration
GO_VERSION ?= 1.21
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build Configuration
BUILD_FLAGS ?= -ldflags="-s -w"
BUILD_OUTPUT ?= $(BINARY_NAME)

# Test Configuration
TEST_FLAGS ?= -v
COVERAGE_OUTPUT ?= coverage.out
COVERAGE_HTML ?= coverage.html

# Linting Configuration
LINT_CONFIG ?= .golangci.yml
LINT_FLAGS ?= run

# Migration Configuration
MIGRATION_EXT ?= sql
MIGRATION_SEQ ?= -seq
MIGRATION_DIGITS ?= 3

# Development Tools
AIR_CONFIG ?= .air.toml
GOLANGCI_LINT_VERSION ?= latest
MIGRATE_VERSION ?= latest
