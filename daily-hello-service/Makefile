SERVICE_NAME ?= ams

start:
	go run cmd/server/main.go

install-deps:
	$ go install github.com/swaggo/swag/cmd/swag@latest
	$ npm install -D semantic-release @semantic-release/gitlab @semantic-release/changelog @semantic-release/git

mod:
	go mod tidy
	go mod vendor

specs: ## Generate swagger specs
	HOST=$(HOST) sh scripts/specs-gen.sh

gen-doc:
	tbls doc --force

# --- Code Actions ------------------------------------------------------------------------
gen-mocks:
	mockgen -source=./internal/repository/cashback.go -destination=./internal/repository/mocks/cashback.go

gen:
	mockgen --version
	go generate -v ./...

lint: 
	$$(go env GOPATH)/bin/golangci-lint version
	$$(go env GOPATH)/bin/golangci-lint run -v -c golangci.yml ./...

gen-specs:
	$$(go env GOPATH)/bin/swag init -g cmd/server/main.go --pd -o swagger-docs

test:
	go test ./... -count=1 -cover -json | $$(go env GOPATH)/bin/tparse -all

test-coverage:
#-- test-coverage for linux --
	go test ./... -count=1 -cover -json -coverprofile=./test_summary/cover.out | $$(go env GOPATH)/bin/tparse -all
	go tool cover -html=./test_summary/cover.out -o ./test_summary/cover.html

test-coverage-windows:
#-- test-coverage for windows --
	go test ./... -count=1 -cover -json -coverprofile="./test_summary/cover.out" | $$(go env GOPATH)/bin/tparse -all
	go tool cover -html="./test_summary/cover.out" -o "./test_summary/cover.html"

#-- Helm --
uat.ams.helm-install:
	helm install ams-api-uat ./helm -f ./helm/ams/uat/values.yaml

uat.consumer.helm-install:
	helm install ams-consumer-uat ./helm -f ./helm/ams-consumer/uat/values.yaml

uat.ams.helm-upgrade:
	@sh ./deploy-helm-manual.sh ams-api ams uat

uat.consumer.helm-upgrade:
	@sh ./deploy-helm-manual.sh ams-consumer ams-consumer uat

clean: ## Clean up the built & test files
	rm -rf ./bootstrap functions/trigger/s3/bootstrap
	rm -rf functions/tfdeployment/dev/*.zip functions/tfdeployment/uat/*.zip functions/tfdeployment/prod/*.zip

govulncheck:
	govulncheck ./...