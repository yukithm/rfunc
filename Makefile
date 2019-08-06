PREFIX = /usr/local
BINDIR = $(PREFIX)/bin

GO_FLAGS = -ldflags="-s -w"

DEVTOOLS_DIR = $(CURDIR)/devtools
DEVTOOLS_BIN = $(DEVTOOLS_DIR)/bin

GOX = $(DEVTOOLS_BIN)/gox
OSARCH = linux/amd64 linux/arm darwin/amd64 windows/386 windows/amd64
DIST_FORMAT = $(DISTDIR)/{{.Dir}}-{{.OS}}-{{.Arch}}
DISTDIR = releases

.PHONY: build
build: rfunc

rfunc: *.go */*.go
	go build $(GO_FLAGS)

.PHONY: install
install: build
	install -d $(BINDIR)
	install rfunc $(BINDIR)

.PHONY: clean
clean:
	rm -f rfunc

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: dist
dist: $(GOX)
	$(GOX) -osarch="$(OSARCH)" $(GO_FLAGS) -output="$(DIST_FORMAT)" .

.PHONY: dist-clean
dist-clean:
	rm -rf rfunc rfunc.exe $(DISTDIR) $(DEVTOOLS_BIN)

.PHONY: proto
proto:
	PATH=$(DEVTOOLS_BIN):$(PATH) protoc -I rfuncs/ rfuncs/rfuncs.proto --go_out=plugins=grpc:rfuncs

.PHONY: mock
mock: proto
	install -d mock_rfuncs
	$(DEVTOOLS_BIN)/mockgen -source=rfuncs/rfuncs.pb.go RFuncsClient > mock_rfuncs/rfuncs_mock.go

.PHONY: init-devenv
init-devenv: devtools

$(GOX): devtools

.PHONY: devtools
devtools:
	go generate $(DEVTOOLS_DIR)/devtools.go
