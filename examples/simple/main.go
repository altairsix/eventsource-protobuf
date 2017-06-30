package main

import (
	"encoding/json"
	"os"
)

//go:generate protoc --go_out=. events.proto
//go:generate protoc --eventsource_out=. events.proto

func main() {
	event := &A{Id: "abc"}
	serializer := NewSerializer()
	record, _ := serializer.MarshalEvent(event)
	actual, _ := serializer.UnmarshalEvent(record)
	json.NewEncoder(os.Stderr).Encode(event)
	json.NewEncoder(os.Stderr).Encode(actual)
}
