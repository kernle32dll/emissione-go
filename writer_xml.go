package emissione

import (
	"encoding/xml"
	"io"
)

// NewXmlWriter creates a new SimpleWriter, marshalling via xml.Marshal.
func NewXmlWriter(options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/xml;charset=utf-8"),
		MarshallMethod(xml.Marshal),

		StreamMethod(func(v interface{}) (io.Reader, error) {
			piper, pipew := io.Pipe()
			enc := xml.NewEncoder(pipew)

			// Start filling the pipe
			go fillXMLPipe(pipew, enc, v)

			return piper, nil
		}),
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

		StreamMethod(func(v interface{}) (io.Reader, error) {
			piper, pipew := io.Pipe()
			enc := xml.NewEncoder(pipew)
			enc.Indent(prefix, indent)

			// Start filling the pipe
			go fillXMLPipe(pipew, enc, v)

			return piper, nil
		}),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}

func fillXMLPipe(pipew *io.PipeWriter, enc *xml.Encoder, v interface{}) {
	defer func() {
		if err := pipew.Close(); err != nil {
			panic(err)
		}
	}()

	if err := enc.Encode(v); err != nil {
		if err := pipew.CloseWithError(err); err != nil {
			panic(err)
		}
	}
}
