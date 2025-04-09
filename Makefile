SHELL := /bin/bash
.PHONY: help config dbcode build migration-up migration-down run test alpha beta minor patch release

help:
	@echo "Makefile Commands:"
	@echo "  config               - Set up the environment."
	@echo "  dbcode               - Generate database code from sql/queries."
	@echo "  migration-up         - Run up migration"
	@echo "  migration-down       - Run down migration"
	@echo "  build                - Build go-rssagg"	
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

dbcode:
	@echo "Generating database code from sql/queries"
	~/go/bin/sqlc generate

migration-up:
	@echo "Running up migration"
	pushd sql/schema && source ../../.env && ~/go/bin/goose postgres $$DB_URL up && popd

migration-down:
	@echo "Running down migration"
	pushd sql/schema && source ../../.env && ~/go/bin/goose postgres $$DB_URL down && popd

build:
	@echo "Building go-rssagg"
	go build -C cmd/go-rssagg

run:
	@echo "Running go-rssagg"
	go build -C cmd/go-rssagg && cmd/go-rssagg/go-rssagg --env .env
	
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
