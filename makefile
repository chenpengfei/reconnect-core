# https://github.com/hashicorp/terraform/blob/master/Makefile

VERSION?=v0.0.1
TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
TESTARGS?=-gcflags=-l

default: test

init:
	@sh -c "'$(CURDIR)/scripts/init.sh'"

tools:
	go get github.com/mattn/goveralls
	go get github.com/golang/mock/mockgen
	go get golang.org/x/tools/cmd/cover
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $$(go env GOPATH)/bin v1.21.0

download:
	go mod tidy
	go mod download

run:
	go run cmd/app/main.go

benchmark:
	go test -bench=. pkg/tools/cal_test.go pkg/tools/cal.go

lint:
	golangci-lint run ./...

test:
	go test -v -race $(TEST) $(TESTARGS) -covermode=atomic

coverprofile:
	go test -v -race $(TEST) $(TESTARGS) -covermode=atomic -coverprofile=coverage.txt

coverage: coverprofile
	go tool cover -html=coverage.txt
	rm coverage.txt

report-coverage: coverprofile
	goveralls -coverprofile=coverage.txt -service=travis-ci
	rm coverage.txt

semantic-release:
	npx semantic-release --no-ci

fmt:
	gofmt -w $(GOFMT_FILES)

# generate runs `go generate` to build the dynamically generated
# source files, except the protobuf stubs, which are built instead with
# "make protobuf".
generate:
	go generate ./...

# disallow any parallelism (-j) for Make. This is necessary since some
# commands during the build process create temporary files that collide
# under parallel conditions.
.NOTPARALLEL:

.PHONY: default

