# eventsource-protobuf

protoc plugin to generate eventsource.Serializers for github.com/altairsix/eventsource

## Installation

```
go get github.com/altairsix/eventsource-protobuf/...
```

## Usage
 
The simplest usage is to use eventsource-protobuf is along side your protoc go:generate line.
  
```text
//go:generate protoc --go_out=. events.proto
//go:generate protoc --plugin=protoc-gen-custom=$GOPATH/bin/eventsource-protobuf --custom_out=. events.proto
```

### Limitations

* All events must be defined as protobuf messages
* All events must be defined in a single .proto file
* One message must be the container, effectively a union type with the following:
    * ```int32 type = 1;```
    * Each event must be included in this message
* Each event must the following fields defined:
    * ```string id``` - aggregate id 
    * ```int32 version``` - version  
    * ```int64 at``` - unix epoch in seconds
    
### Example

The following provides a simple example of a .proto that adheres to all the stated
requirements.

```proto
message ItemAdded {
    string id = 1;
    int32 version = 2;
    int64 at = 3;
}

message ItemRemoved {
    string id = 1;
    int32 version = 2;
    int64 at = 3;
}

// Events 
message Events {
    int32 type = 1;
    ItemAdded a = 2;
    ItemRemoved b = 3;
}
```