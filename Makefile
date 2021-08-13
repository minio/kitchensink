PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)

GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

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

# Builds kitchensink locally.
build: checks
	@echo "Building kitchensink binary to './kitchensink'"
	@GO111MODULE=on CGO_ENABLED=0 go build -trimpath -tags kqueue --ldflags "-s -w" -o $(PWD)/kitchensink

# Builds kitchensink and installs it to $GOPATH/bin.
install: build
	@echo "Installing kitchensink binary to '$(GOPATH)/bin/kitchensink'"
	@mkdir -p $(GOPATH)/bin && cp -f $(PWD)/kitchensink $(GOPATH)/bin/kitchensink
	@echo "Installation successful. To learn more, try \"kitchensink --help\"."

clean:
	@echo "Cleaning up all the generated files"
	@find . -name '*.txt' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@rm -rvf kitchensink
	@rm -rvf build
	@rm -rvf release
