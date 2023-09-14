all: run

build:
	@docker compose build

down:
	@docker compose down
	
run: build
	@docker compose up -d

test:
	@go test -v -coverpkg ./... -coverprofile=profile.cov ./...
	@cat profile.cov | grep -v mocks > profile.filtred.cov
	
lint:
	@golangci-lint run ./...

