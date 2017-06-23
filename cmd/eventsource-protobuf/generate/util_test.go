package generate

import (
	"testing"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/stretchr/testify/assert"
)

func TestCamel(t *testing.T) {
	testCases := map[string]struct {
		In       string
		Expected string
	}{
		"simple": {
			In:       "hello",
			Expected: "Hello",
		},
		"compound": {
			In:       "hello_world",
			Expected: "HelloWorld",
		},
		"double _": {
			In:       "hello__world",
			Expected: "HelloWorld",
		},
		"many": {
			In:       "a_b_c",
			Expected: "ABC",
		},
	}

	for label, tc := range testCases {
		t.Run(label, func(t *testing.T) {
			assert.Equal(t, tc.Expected, camel(tc.In))
		})
	}
}

func TestIDFields(t *testing.T) {
	results := idFields(&descriptor.FileDescriptorProto{
		MessageType: []*descriptor.DescriptorProto{
			{
				Name: String("Foo"),
				Field: []*descriptor.FieldDescriptorProto{
					{
						Name: String("ID"),
					},
				},
			},
		},
	})
	assert.Equal(t, map[string]string{"Foo": "ID"}, results)
}
