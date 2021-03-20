
.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: build
## build: builds server
build:
	@cd ./cmd/server; \
	go build .

.PHONY: build-static
## build-static: build for running binary insite scratch container (runned by DroneCI)
build-static:
	@cd ./cmd/server; \
	CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o server .

.PHONY: run
## run: run locally (don't forget to set all required environment variables)
run:
	@cd ./cmd/server; \
	go run .

.PHONY: vet
## vet: runs go vet command
vet:
	@go vet ./cmd/...

.PHONY: test
## test: runs go vet and go test commands
test: vet
	@go test ./...  -coverprofile cp.out

.PHONY: build-docker
## build-docker: builds Docker image
build-docker:
	@docker build --tag gbfs-tools:latest ./cmd/server;
