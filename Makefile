SHELL := /bin/bash
.PHONY: help config build migration-up migration-down run test alpha beta minor patch release

help:
	@echo "Makefile Commands:"
	@echo "  config               - Set up the environment."
	@echo "  build                - Build go-rssagg"
	@echo "  migration-up         - Run up migration"
	@echo "  migration-down       - Run down migration"
	@echo "  run                  - Build and run go-rssagg in-place"
	@echo "  test                 - Run tests"
	@echo "  alpha                - Generate changelog and create an alpha tag."
	@echo "  beta                 - Generate changelog and create an beta tag."
	@echo "  minor                - Generate changelog and create a minor tag."
	@echo "  patch                - Generate changelog and create a patch tag."
	@echo "  release              - Generate changelog and create a release tag."

all: config build run

config:
	@echo "Installing required tools"
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go mod tidy
	go mod vendor

build:
	@echo "Building go-rssagg"
	go build

migration-up:
	@echo "Running up migration"
	pushd sql/schema && source ../../.env && ~/go/bin/goose postgres $$DB_URL up && popd

migration-down:
	@echo "Running down migration"
	pushd sql/schema && source ../../.env && ~/go/bin/goose postgres $$DB_URL down && popd

run:
	@echo "Running go-rssagg"
	go build && ./go-rssagg
	
test:
	@echo "Running tests"
#	./util/test.sh

alpha:
	@echo "Generating changelog and tag"
	commit-and-tag-version --prerelease alpha

beta:
	@echo "Generating changelog and tag"
	commit-and-tag-version --prerelease beta

minor:
	@echo "Generating changelog and tag"
	commit-and-tag-version --release-as minor

patch:
	@echo "Generating changelog and tag"
	commit-and-tag-version --release-as patch

release:
	@echo "Generating changelog and tag"
	commit-and-tag-version
