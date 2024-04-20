include .env

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

confirm: 
	@echo 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

##race: run the cmd/api application with race flag
race:
	go run -race ./cmd/api -db-dsn=${db_dsn}

## run/api: run the cmd/api application
run: confirm
	go run ./cmd/api -db-dsn=${db_dsn}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code'
	go fmt ./...
	@echo 'Vetting code'
	go vet ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...