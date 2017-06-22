package main

import (
	"fmt"
	"time"

	"github.com/altairsix/eventsource"
	"github.com/gogo/protobuf/proto"
)

type serializer struct {
}

func (s *serializer) MarshalEvent(event eventsource.Event) (eventsource.Record, error) {
	container := &Events{};

	switch v := event.(type) {

	case *A:
		container.Type = 2
		container.A = v

	case *B:
		container.Type = 3
		container.B = v

	default:
		return eventsource.Record{}, fmt.Errorf("Unhandled type, %v", event)
	}

	data, err := proto.Marshal(container)
	if err != nil {
		return eventsource.Record{}, err
	}

	return eventsource.Record{
		Version: event.EventVersion(),
		Data:    data,
	}, nil
}

func (s *serializer) UnmarshalEvent(record eventsource.Record) (eventsource.Event, error) {
	container := &Events{};
	err := proto.Unmarshal(record.Data, container)
	if err != nil {
		return nil, err
	}

	var event interface{}
	switch container.Type {

	case 2:
		event = container.A

	case 3:
		event = container.B

	default:
		return nil, fmt.Errorf("Unhandled type, %v", container.Type)
	}

	return event.(eventsource.Event), nil
}

func NewSerializer() eventsource.Serializer {
	return &serializer{}
}

func (e *A) AggregateID() string { return e.Id }
func (e *A) EventVersion() int   { return int(e.Version) }
func (e *A) EventAt() time.Time  { return time.Unix(e.At, 0) }

func (e *B) AggregateID() string { return e.Id }
func (e *B) EventVersion() int   { return int(e.Version) }
func (e *B) EventAt() time.Time  { return time.Unix(e.At, 0) }

