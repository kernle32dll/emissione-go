package emissione

import (
	"encoding/xml"
)

func NewXmlWriter(options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/xml;charset=utf-8"),
		MarshallMethod(xml.Marshal),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}

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
