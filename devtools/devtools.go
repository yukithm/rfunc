// +build devtools

package devtools

import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/mitchellh/gox"
)

//go:generate go build -v -o=./bin/gox github.com/mitchellh/gox
//go:generate go build -v -o=./bin/protoc-gen-go github.com/golang/protobuf/protoc-gen-go
//go:generate go build -v -o=./bin/mockgen github.com/golang/mock/mockgen
