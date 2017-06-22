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

func capitalize(in string) string {
	if in == "" {
		return ""
	}
	return strings.ToUpper(in[0:1]) + in[1:]
}

func newTemplate(content string) (*template.Template, error) {
	fn := map[string]interface{}{
		"base":       base,
		"lower":      lower,
		"capitalize": capitalize,
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
