GOPATH := $(shell go env GOPATH)
GO_FLAGS := -ldflags="-s -w"
REPO := github.com/yukithm/rfunc

PREFIX := /usr/local
BINDIR := $(PREFIX)/bin

DEVTOOL_DIR = $(CURDIR)/devtool
GOX = $(DEVTOOL_DIR)/bin/gox
OSARCH = linux/amd64 linux/arm darwin/amd64 windows/386 windows/amd64
DISTDIR := releases
DIST_FORMAT = $(DISTDIR)/{{.Dir}}-{{.OS}}-{{.Arch}}

.PHONY: build rfunc install clean proto mock test dist dist-clean

build: rfunc

rfunc: *.go */*.go
	go build $(GO_FLAGS)

install: build
	install -d $(BINDIR)
	install rfunc $(BINDIR)

clean:
	rm -f rfunc

proto:
	go get github.com/golang/protobuf/protoc-gen-go
	protoc -I rfuncs/ rfuncs/rfuncs.proto --go_out=plugins=grpc:rfuncs

mock:
	go get github.com/golang/mock/gomock
	go get github.com/golang/mock/mockgen
	install -d mock_rfuncs
	mockgen $(REPO)/rfuncs RFuncsClient > mock_rfuncs/rfuncs_mock.go

test:
	go test -v ./... -cover

dist: $(DEVTOOL_DIR)/bin/gox
	$(GOX) -osarch="$(OSARCH)" $(GO_FLAGS) -output="$(DIST_FORMAT)" .

$(DEVTOOL_DIR)/bin/gox:
	mkdir -p $(DEVTOOL_DIR)/{bin,pkg,src}
	GOPATH=$(DEVTOOL_DIR) go get github.com/mitchellh/gox

dist-clean:
	rm -rf rfunc $(DISTDIR) $(DEVTOOL_DIR)
