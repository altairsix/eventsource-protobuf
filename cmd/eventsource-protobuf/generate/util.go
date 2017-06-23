package generate

import (
	"errors"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
)

// String returns a pointer to the provide string or nil if the zero value was passed in
func String(in string) *string {
	if in == "" {
		return nil
	}

	return &in
}

func name(in *descriptor.FileDescriptorProto) *string {
	if in.Name != nil {
		name := *in.Name
		ext := filepath.Ext(name)
		name = name[0 : len(name)-len(ext)]
		return String(name + ".pb.es.go")
	}

	return String("events.pb.es.go")
}

func packageName(in *descriptor.FileDescriptorProto) (string, error) {
	if in.Package != nil {
		return *in.Package, nil
	}

	if in.Name != nil {
		name := *in.Name
		ext := filepath.Ext(name)
		return name[0 : len(name)-len(ext)], nil
	}

	return "", errors.New("unable to determine package name")
}

func base(in string) string {
	idx := strings.LastIndex(in, ".")
	if idx == -1 {
		return in
	}
	return in[idx+1:]
}

func lower(in string) string {
	return strings.ToLower(in)
}

func camel(in string) string {
	segments := strings.Split(in, "_")
	capped := make([]string, 0, len(segments))

	for _, segment := range segments {
		if segment == "" {
			continue
		}
		capped = append(capped, strings.ToUpper(segment[0:1])+segment[1:])
	}
	return strings.Join(capped, "")
}

func newTemplate(content string) (*template.Template, error) {
	fn := map[string]interface{}{
		"base":  base,
		"lower": lower,
		"camel": camel,
	}

	return template.New("page").Funcs(fn).Parse(content)
}

// findContainerMessage returns the message that contains all the other message types
func findContainerMessage(in *descriptor.FileDescriptorProto) (*descriptor.DescriptorProto, error) {
outer:
	for _, message := range in.MessageType {
		for index, field := range message.Field {
			if index > 0 {
				return nil, errors.New("not found")
			}
			if *field.Name != "type" || *field.Number != int32(1) {
				continue outer
			}
			return message, nil
		}
	}

	return nil, errors.New("Not found")
}

// idFields returns a map of Type -> ID field name
func idFields(in *descriptor.FileDescriptorProto) map[string]string {
	results := map[string]string{}

outer:
	for _, mType := range in.MessageType {
		for _, field := range mType.Field {
			if strings.ToLower(*field.Name) == "id" {
				results[*mType.Name] = *field.Name
				continue outer
			}
		}
	}

	return results
}
