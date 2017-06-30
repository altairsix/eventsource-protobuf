package main

import (
	"encoding/json"
	"os"
)

//go:generate protoc -I .:$GOPATH/src:$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=. event.proto
//go:generate protoc -I .:$GOPATH/src:$GOPATH/src/github.com/gogo/protobuf/protobuf --eventsource_out=. event.proto

func main() {
	event := &A{ID: "abc"}
	serializer := NewSerializer()
	record, _ := serializer.MarshalEvent(event)
	actual, _ := serializer.UnmarshalEvent(record)
	json.NewEncoder(os.Stderr).Encode(event)
	json.NewEncoder(os.Stderr).Encode(actual)
}
