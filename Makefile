# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

# ==================================================================================== #
# BUILD
# ==================================================================================== #

git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.version=${git_description}'

.PHONY: build/api
build/api:
	@echo "Building cmd/api..."
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o="./bin/linux_amd64/api" "./cmd/api"
	GOOS=darwin GOARCH=amd64 go build -ldflags=${linker_flags} -o="./bin/osx_amd64/api" "./cmd/api"
	GOOS=windows GOARCH=amd64 go build -ldflags=${linker_flags} -o="./bin/windows_amd64/api" "./cmd/api"

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
.PHONY: audit
audit:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...