package emissione

import (
	"encoding/xml"
)

// NewXmlWriter creates a new SimpleWriter, marshalling via xml.Marshal.
func NewXmlWriter(options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/xml;charset=utf-8"),
		MarshallMethod(xml.Marshal),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}

// NewXmlIndentWriter creates a new SimpleWriter, marshalling via xml.MarshalIndent.
func NewXmlIndentWriter(prefix string, indent string, options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/xml;charset=utf-8"),

		// Define a custom marshall method, which uses MarshalIndent
		MarshallMethod(func(v interface{}) ([]byte, error) {
			return xml.MarshalIndent(v, prefix, indent)
		}),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}
