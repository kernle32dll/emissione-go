package emissione

import (
	"encoding/json"
	"io"
)

// NewJSONWriter creates a new SimpleWriter, marshalling via json.Marshal.
func NewJSONWriter(options ...WriterOption) Writer {
	defaults := []WriterOption{
		ContentType("application/json;charset=utf-8"),
		MarshallMethod(json.Marshal),
		StreamMethod(func(v interface{}) (io.Reader, error) {
			piper, pipew := io.Pipe()
			enc := json.NewEncoder(pipew)

			// Start filling the pipe
			go fillJSONPipe(pipew, enc, v)

			return piper, nil
		}),
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

		StreamMethod(func(v interface{}) (io.Reader, error) {
			piper, pipew := io.Pipe()
			enc := json.NewEncoder(pipew)
			enc.SetIndent(prefix, indent)

			// Start filling the pipe
			go fillJSONPipe(pipew, enc, v)

			return piper, nil
		}),
	}

	return NewSimpleWriter(append(defaults, options...)...)
}

func fillJSONPipe(pipew *io.PipeWriter, enc *json.Encoder, v interface{}) {
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
