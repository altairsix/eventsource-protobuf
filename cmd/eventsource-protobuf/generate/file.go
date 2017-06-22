package generate

import (
	"bytes"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/plugin"
)

const (
	content = `package {{ .Package }}

import (
	"fmt"
	"time"

	"github.com/altairsix/eventsource"
	"github.com/gogo/protobuf/proto"
)

type serializer struct {
}

func (s *serializer) MarshalEvent(event eventsource.Event) (eventsource.Record, error) {
	container := &{{ .Message.Name | base | camel }}{}

	switch v := event.(type) {
{{ range .Fields }}
	case *{{ .TypeName | base | camel }}:
		container.Type = {{ .Number }}
		container.{{ .Name | camel }} = v
{{ end }}
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
	container := &{{ .Message.Name | base | camel }}{};
	err := proto.Unmarshal(record.Data, container)
	if err != nil {
		return nil, err
	}

	var event interface{}
	switch container.Type {
{{ range .Fields }}
	case {{ .Number }}:
		event = container.{{ .Name | camel }}
{{ end }}
	default:
		return nil, fmt.Errorf("Unhandled type, %v", container.Type)
	}

	return event.(eventsource.Event), nil
}

func NewSerializer() eventsource.Serializer {
	return &serializer{}
}
{{ range .Fields }}
func (m *{{ .TypeName | base | camel }}) AggregateID() string { return m.Id }
func (m *{{ .TypeName | base | camel }}) EventVersion() int   { return int(m.Version) }
func (m *{{ .TypeName | base | camel }}) EventAt() time.Time  { return time.Unix(m.At, 0) }
{{ end }}
`
)

// File accepts the proto file definition and returns the response for this file
func File(in *descriptor.FileDescriptorProto) (*plugin_go.CodeGeneratorResponse_File, error) {
	pkg, err := packageName(in)
	if err != nil {
		return nil, err
	}

	message, err := findContainerMessage(in)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	t, err := newTemplate(content)
	if err != nil {
		return nil, err
	}

	t.Execute(buf, map[string]interface{}{
		"Package": pkg,
		"Message": message,
		"Fields":  message.Field[1:],
	})

	return &plugin_go.CodeGeneratorResponse_File{
		Name:    name(in),
		Content: String(buf.String()),
	}, nil
}

// AllFiles accepts multiple proto file definitions and returns the list of files
func AllFiles(in []*descriptor.FileDescriptorProto) ([]*plugin_go.CodeGeneratorResponse_File, error) {
	results := make([]*plugin_go.CodeGeneratorResponse_File, 0, len(in))

	if in != nil {
		for _, file := range in {
			v, err := File(file)
			if err != nil {
				return nil, err
			}
			results = append(results, v)
		}
	}

	return results, nil
}
