NAME=dig-up
BINDIR=dist

PACKAGE=github.com/Ehco1996/ehco/internal/constant
BUILDTIME=$(shell date +"%m-%d-%y-%T")
BRANCH=$(shell git rev-parse --abbrev-ref HEAD | tr -d '\040\011\012\015\n')
REVISION=$(shell git rev-parse --short HEAD)

# -w -s 参数的解释：You will get the smallest binaries if you compile with -ldflags '-w -s'. The -w turns off DWARF debugging information
# for more information, please refer to https://stackoverflow.com/questions/22267189/what-does-the-w-flag-mean-when-passed-in-via-the-ldflags-option-to-the-go-comman
GOBUILD=CGO_ENABLED=0 go build -trimpath -ldflags="-w -s -X ${PACKAGE}.GitBranch=${BRANCH} -X ${PACKAGE}.GitRevision=${REVISION} -X ${PACKAGE}.BuildTime=${BUILDTIME}"

.PHONY: fmt test build tidy

fmt:
	golangci-lint run --fix

test:
	go test -v -count=1  -coverpkg=./internal -timeout=10s ./...

build:
	${GOBUILD} -o $(BINDIR)/$(NAME) cmd/dig-up/main.go

build-linux-amd64:
	GOARCH=amd64 GOOS=linux ${GOBUILD} -o $(BINDIR)/$(NAME)_amd64 cmd/dig-up/main.go

tidy:
	go mod tidy
