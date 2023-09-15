all: run update-cert

build:
	@docker compose build

down: clean
	@docker compose down
	
create-crt:
	@mkdir certs
	@touch ./certs/ca.crt
run:clean create-crt build 
	@docker compose up -d

test:
	@go test -v -coverpkg ./... -coverprofile=profile.cov ./...
	@cat profile.cov | grep -v mocks > profile.filtred.cov
	
clean: 
	@rm -rf certs
	
update-cert:
	@sudo cp certs/ca.crt /etc/pki/ca-trust/source/anchors/proxyca.pem
	@sudo update-ca-trust
	
lint:
	@golangci-lint run ./...

