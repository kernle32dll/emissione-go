package emissione

import (
	"fmt"
	"net/http"
)

type SimpleWriter struct {
	marshallMethod func(v interface{}) ([]byte, error)
	contentType    string
}

func (writer SimpleWriter) Write(w http.ResponseWriter, i interface{}) error {
	// Set content type, if not already set
	if len(w.Header().Get("Content-Type")) == 0 {
		w.Header().Set("Content-Type", writer.contentType)
	}

	bytes, err := writer.marshallMethod(i)
	if err != nil {
		return err
	}

	_, wErr := w.Write(bytes)
	return wErr
}

func NewSimpleWriter(setters ...WriterOption) Writer {
	// Default Options
	args := &WriterOptions{
		MarshallMethod: func(v interface{}) (bytes []byte, e error) {
			return []byte(fmt.Sprint(v)), nil
		},
		ContentType: "application/json",
	}

	for _, setter := range setters {
		setter(args)
	}

	return &SimpleWriter{
		marshallMethod: args.MarshallMethod,
		contentType:    args.ContentType,
	}
}
