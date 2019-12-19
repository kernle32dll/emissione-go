package emissione

import (
	"fmt"
	"io"
	"net/http"
)

// SimpleWriter is a simple implementation of emissiones Writer interface,
// delegating writing duty to a specific marshaller, and setting the appropriate
// content type header.
type SimpleWriter struct {
	marshallMethod func(v interface{}) ([]byte, error)
	streamMethod   func(v interface{}) (io.Reader, error)
	contentType    string
}

func (writer SimpleWriter) Write(w http.ResponseWriter, i interface{}) (returnErr error) {
	// Set content type, if not already set
	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", writer.contentType)
	}

	if writer.streamMethod != nil {
		stream, err := writer.streamMethod(i)
		if err != nil {
			return err
		}

		if autoClose, isAutoClose := stream.(io.ReadCloser); isAutoClose {
			defer func() {
				returnErr = autoClose.Close()
			}()
		}

		_, sErr := io.Copy(w, stream)
		return sErr
	} else {
		bytes, err := writer.marshallMethod(i)
		if err != nil {
			return err
		}

		_, wErr := w.Write(bytes)
		return wErr
	}
}

// NewSimpleWriter instantiates a new SimpleWriter.
//
// See the documentation for WriterOption for configuration options.
func NewSimpleWriter(setters ...WriterOption) Writer {
	// Default Options
	args := &WriterOptions{
		MarshallMethod: func(v interface{}) (bytes []byte, e error) {
			return []byte(fmt.Sprint(v)), nil
		},
		StreamMethod: nil,
		ContentType:  "text/plain",
	}

	for _, setter := range setters {
		setter(args)
	}

	return &SimpleWriter{
		marshallMethod: args.MarshallMethod,
		streamMethod:   args.StreamMethod,
		contentType:    args.ContentType,
	}
}
