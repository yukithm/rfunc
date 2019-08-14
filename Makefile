NAME := rfunc
PREFIX := /usr/local
BINDIR := $(PREFIX)/bin

VERSION := $(shell git describe --tags --always --dirty=-dev)
LDFLAGS := -s -w -X 'github.com/yukithm/rfunc/cmd.version=$(VERSION)'

DEVTOOLS_DIR := $(CURDIR)/devtools
DEVTOOLS_BIN := $(DEVTOOLS_DIR)/bin

DISTDIR := releases
OSARCH := linux/amd64 linux/arm darwin/amd64 windows/386 windows/amd64
DIST_FORMAT := $(DISTDIR)/{{.Dir}}-{{.OS}}-{{.Arch}}

SOURCES := $(shell find . -type f -name "*.go")

export GO111MODULE=on

.PHONY: build
build: $(NAME)

$(NAME): $(SOURCES)
	go build -ldflags "$(LDFLAGS)"

.PHONY: install
install: build
	install -d $(BINDIR)
	install $(NAME) $(BINDIR)

.PHONY: dist
dist: devtools
	$(DEVTOOLS_BIN)/gox -osarch="$(OSARCH)" -ldflags "$(LDFLAGS)" -output="$(DIST_FORMAT)" .

.PHONY: dist-clean
dist-clean:
	rm -rf $(NAME) $(NAME).exe "$(DISTDIR)" "$(DEVTOOLS_BIN)"

.PHONY: clean
clean:
	rm -f $(NAME) $(NAME).exe

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: deps
deps: download-deps devtools

.PHONY: download-deps
download-deps:
	go get -v -d

.PHONY: devtools
devtools:
	go generate $(DEVTOOLS_DIR)/devtools.go

.PHONY: lint
lint:
	go vet ./...
	golint -set_exit_status ./...

.PHONY: proto
proto:
	PATH=$(DEVTOOLS_BIN):$(PATH) protoc -I rfuncs/ rfuncs/rfuncs.proto --go_out=plugins=grpc:rfuncs

.PHONY: mock
mock: proto
	install -d mock_rfuncs
	$(DEVTOOLS_BIN)/mockgen -source=rfuncs/rfuncs.pb.go RFuncsClient > mock_rfuncs/rfuncs_mock.go
