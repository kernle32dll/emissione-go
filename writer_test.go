package emissione_test

import (
	"errors"
	"fmt"
	"github.com/kernle32dll/emissione-go"
	"net/http/httptest"
	"testing"
)

func TestSimpleWriter_Write(t *testing.T) {
	type SampleObject struct {
		TestProp int `json:"testProp",xml:"testProp"`
	}

	sampleObject := SampleObject{42}
	simpleMarshalled := fmt.Sprint(sampleObject)

	t.Run("simple", func(t *testing.T) {
		w := httptest.NewRecorder()

		writer := emissione.NewSimpleWriter()

		if err := writer.Write(w, sampleObject); err != nil {
			t.Errorf("Write() error = %v, wantErr %v", err, nil)
		}

		if got := w.Body.String(); got != simpleMarshalled {
			t.Errorf("Write() = %s, want %s", got, simpleMarshalled)
		}

		if got := w.Header().Get("Content-Type"); got != "text/plain" {
			t.Errorf("Write() = %s, want %s", got, "text/plain")
		}
	})

	t.Run("existing-content-type", func(t *testing.T) {
		w := httptest.NewRecorder()
		w.Header().Add("Content-Type", "text/plain")

		writer := emissione.NewSimpleWriter(
			emissione.ContentType("not-used"),
		)

		if err := writer.Write(w, sampleObject); err != nil {
			t.Errorf("Write() error = %v, wantErr %v", err, nil)
		}

		if got := w.Body.String(); got != simpleMarshalled {
			t.Errorf("Write() = %s, want %s", got, simpleMarshalled)
		}

		if got := w.Header().Get("Content-Type"); got != "text/plain" {
			t.Errorf("Write() = %s, want %s", got, "text/plain")
		}
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()

		expectedErr := errors.New("error")

		writer := emissione.NewSimpleWriter(
			emissione.MarshallMethod(func(_ interface{}) ([]byte, error) {
				return []byte("content-to-be-ignored"), expectedErr
			}),
		)

		if err := writer.Write(w, sampleObject); err != expectedErr {
			t.Errorf("Write() error = %v, wantErr %v", err, expectedErr)
		}

		if got := w.Body.String(); got != "" {
			t.Errorf("Write() = %s, want %s", got, "")
		}

		if got := w.Header().Get("Content-Type"); got != "text/plain" {
			t.Errorf("Write() = %s, want %s", got, "text/plain")
		}
	})
}
