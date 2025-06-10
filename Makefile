GO ?= $(shell command -v go 2> /dev/null)
PACKAGES=$(shell go list ./...)

## Checks the code style and tests
all: check-style test

## Runs govet and gofmt against all packages.
.PHONY: check-style
check-style: govet goformat
	@echo Checking for style guide compliance

## Runs govet against all packages.
.PHONY: vet
govet:
	@echo Running govet
	$(GO) vet ./...
	@echo Govet success

## Checks if files are formatted with go fmt.
.PHONY: goformat
goformat:
	@echo Running gofmt
	@for package in $(PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failed"; \
				echo "To fix it, run:"; \
				echo "go fmt [FAILED_PACKAGE]"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \

## Runs go test against all packages.
.PHONY: test
test:
	@echo Running go tests
	go test ./...
	@echo Go test success