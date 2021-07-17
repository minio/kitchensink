PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

BUILD_LDFLAGS := '$(LDFLAGS)'

all: build

checks:
	@echo "Checking dependencies"
	@(env bash $(PWD)/buildscripts/checkdeps.sh)

getdeps:
	@mkdir -p ${GOPATH}/bin
	@which golangci-lint 1>/dev/null || (echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.27.0)
	@which stringer 1>/dev/null || (echo "Installing stringer" && go get golang.org/x/tools/cmd/stringer)

crosscompile:
	@(env bash $(PWD)/buildscripts/cross-compile.sh)

verifiers: getdeps vet fmt lint

fmt:
	@echo "Running $@"
	@GO111MODULE=on gofmt -d cmd/
	@GO111MODULE=on gofmt -d pkg/

lint:
	@echo "Running $@ check"
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint cache clean
	@GO111MODULE=on ${GOPATH}/bin/golangci-lint run --timeout=5m --config ./.golangci.yml

# Builds ks locally.
build: checks
	@echo "Building ks binary to './ks'"
	@GO111MODULE=on CGO_ENABLED=0 go build -trimpath -tags kqueue --ldflags $(BUILD_LDFLAGS) -o $(PWD)/ks

# Builds KS and installs it to $GOPATH/bin.
install: build
	@echo "Installing ks binary to '$(GOPATH)/bin/ks'"
	@mkdir -p $(GOPATH)/bin && cp -f $(PWD)/ks $(GOPATH)/bin/ks
	@echo "Installation successful. To learn more, try \"ks --help\"."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf ks
	@rm -rvf build
	@rm -rvf release