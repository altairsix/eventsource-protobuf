package main

import (
	"fmt"

	"github.com/altairsix/eventsource"
)

//go:generate protoc -I .:$GOPATH/src:$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. commands.proto
//go:generate protoc -I .:$GOPATH/src:$GOPATH/src/github.com/gogo/protobuf/protobuf --commands_out=. commands.proto

func main() {
	var cmd eventsource.Command

	cmd = &RegisterUser{ID: "abc"}
	fmt.Printf("%#v\n", cmd)

	cmd = &ResetPassword{ID: "def"}
	fmt.Printf("%#v\n", cmd)
}
