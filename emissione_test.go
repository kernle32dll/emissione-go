package emissione_test

import (
	"github.com/kernle32dll/emissione-go"

	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockWriter struct {
	writeOutput string
	errorOut    error
}

func (mockWriter MockWriter) Write(w http.ResponseWriter, i interface{}) error {
	if mockWriter.writeOutput != "" {
		_, err := w.Write([]byte(mockWriter.writeOutput))
		return err
	}

	return mockWriter.errorOut
}

type ErrorResponseWriter struct {
	*httptest.ResponseRecorder
	err error
}

func (w *ErrorResponseWriter) Write(b []byte) (int, error) {
	return 0, w.err
}

func TestHandler_Write(t *testing.T) {
	jsonOut, xmlOut := "json", "xml"

	jsonWriter := &MockWriter{jsonOut, nil}
	xmlWriter := &MockWriter{xmlOut, nil}

	writer := emissione.New(
		jsonWriter,
		emissione.WriterMapping{
			"application/json": jsonWriter,
			"application/xml":  xmlWriter,
		},
	)

	tests := []struct {
		name    string
		accepts []string
		want    string
	}{
		{"default", []string{"unknown"}, jsonOut},
		{"empty-header", nil, jsonOut},
		{"empty-but-different", []string{""}, jsonOut},
		{"*/*", []string{"*/*"}, jsonOut},
		{"application/json", []string{"application/json"}, jsonOut},
		{"application/*", []string{"application/*"}, jsonOut},
		{"application/xml", []string{"application/xml"}, xmlOut},
		{"weighted", []string{"application/*, application/xml;q=2"}, xmlOut},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()

			writer.Write(w, &http.Request{Header: map[string][]string{"Accept": tt.accepts}}, http.StatusTeapot, struct{}{})

			if w.Code != http.StatusTeapot {
				t.Errorf("Write() = %v, want %v", w.Code, http.StatusTeapot)
			}

			if got := w.Body.String(); got != tt.want {
				t.Errorf("Write() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Write_writer_error(t *testing.T) {
	expectedErr := errors.New("some-error")
	errorWriter := &MockWriter{"", expectedErr}

	writer := emissione.New(
		errorWriter,
		nil,
	)

	w := httptest.NewRecorder()

	// Catch panic
	defer func() {
		if got := w.Body.String(); got != "" {
			t.Errorf("Write() = %v, want %v", got, "")
		}

		if r := recover(); r != nil {
			if r != expectedErr {
				t.Errorf("Write_error() = %v, want %v", r, expectedErr)
			}
		}
	}()
	writer.Write(w, &http.Request{}, http.StatusTeapot, struct{}{})

	t.Errorf("Write_error() = %v, want %v", nil, expectedErr)
}

func TestHandler_Write_buffer_error(t *testing.T) {
	errorWriter := &MockWriter{"discarded", nil}

	writer := emissione.New(
		errorWriter,
		nil,
	)

	expectedErr := errors.New("some-error")

	w := &ErrorResponseWriter{
		httptest.NewRecorder(),
		expectedErr,
	}

	httptest.NewRecorder()

	// Catch panic
	defer func() {
		if got := w.Body.String(); got != "" {
			t.Errorf("Write() = %v, want %v", got, "")
		}

		if r := recover(); r != nil {
			if r != expectedErr {
				t.Errorf("Write_error() = %v, want %v", r, expectedErr)
			}
		}
	}()
	writer.Write(w, &http.Request{}, http.StatusTeapot, struct{}{})

	t.Errorf("Write_error() = %v, want %v", nil, expectedErr)
}

func TestHandler_Write_unsupported(t *testing.T) {
	writer := emissione.New(
		nil,
		nil,
	)

	w := httptest.NewRecorder()

	writer.Write(w, &http.Request{}, http.StatusTeapot, struct{}{})

	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("Write_unsupported() = %v, want %v", w.Code, http.StatusUnsupportedMediaType)
	}

	if got := w.Body.String(); got != "" {
		t.Errorf("Write_unsupported() = %v, want %v", got, "")
	}
}
