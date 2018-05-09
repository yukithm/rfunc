GOPATH := $(shell go env GOPATH)
GO_FLAGS := -ldflags="-s -w"
REPO := github.com/yukithm/rfunc

.PHONY: rfunc proto mock test

rfunc: *.go */*.go
	go build $(GO_FLAGS)

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
