package emissione_test

import (
	"github.com/kernle32dll/emissione-go"

	"errors"
	"fmt"
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

		assertState(t, w, writer, sampleObject, "text/plain", simpleMarshalled, nil)
	})

	t.Run("existing-content-type", func(t *testing.T) {
		w := httptest.NewRecorder()
		w.Header().Add("Content-Type", "text/plain")

		writer := emissione.NewSimpleWriter(
			emissione.ContentType("not-used"),
		)

		assertState(t, w, writer, sampleObject, "text/plain", simpleMarshalled, nil)
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()

		expectedErr := errors.New("error")

		writer := emissione.NewSimpleWriter(
			emissione.MarshallMethod(func(_ interface{}) ([]byte, error) {
				return []byte("content-to-be-ignored"), expectedErr
			}),
		)

		assertState(t, w, writer, sampleObject, "text/plain", "", expectedErr)
	})
}

func assertState(
	t *testing.T, w *httptest.ResponseRecorder, writer emissione.Writer,
	object interface{}, expectedContentType string, expectedOutput string, expectedError error,
) {
	if err := writer.Write(w, object); err != expectedError {
		t.Errorf("Write() error = %v, want %v", err, expectedError)
	}
	if got := w.Body.String(); got != expectedOutput {
		t.Errorf("Write() = %s, want %s", got, expectedOutput)
	}
	if got := w.Header().Get("Content-Type"); got != expectedContentType {
		t.Errorf("Write() = %s, want %s", got, expectedContentType)
	}
}
