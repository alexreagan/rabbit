BINARY="rabbit"
VERSION=1.0.0
BUILD=`date +%FT%T%z`

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`
export CGO_ENABLED=0
export GOOS=linux
# export GOOS=windows
# export GOARCH=amd64
#export CGO_ENABLED=1
#export CC=x86_64-linux-musl-gcc
#export CXX=x86_64-linux-musl-g++
#CGO_LDFLAGS="-static"

all: rabbit

rabbit:
	go build -a

list:
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}

fmt:
	@gofmt -s -w ${GOFILES}

fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
 		echo "Please run 'make fmt' and commit the result:"; \
 		echo "$${diff}"; \
 		exit 1; \
	fi;

install:
	@govendor sync -v

test:
	@go test -cpu=1,2,4 -v ./...

vet:
	@go vet $(VETPACKAGES)

docker:
 	@docker build -t AlexReagan/rabbit:latest .

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: default fmt fmt-check install test vet docker clean rabbit
