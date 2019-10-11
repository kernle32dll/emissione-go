package emissione

import (
	"encoding/json"
)

// NewJSONWriter creates a new SimpleWriter, marshalling via json.Marshal.
func NewJSONWriter(options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/json;charset=utf-8"),
		MarshallMethod(json.Marshal),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}

// NewJSONIndentWriter creates a new SimpleWriter, marshalling via json.MarshalIndent.
func NewJSONIndentWriter(prefix string, indent string, options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/json;charset=utf-8"),

		// Define a custom marshall method, which uses MarshalIndent
		MarshallMethod(func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, prefix, indent)
		}),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}
