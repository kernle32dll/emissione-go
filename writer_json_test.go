package emissione_test

import (
	"github.com/kernle32dll/emissione-go"

	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSONWriter(t *testing.T) {
	type SampleObject struct {
		TestProp int `json:"testProp"`
	}

	sampleObject := SampleObject{42}

	t.Run("simple", func(t *testing.T) {
		jsonMarshalled, mErr := json.Marshal(sampleObject)
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewJSONWriter(emissione.StreamMethod(nil))

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestJSONWriter() unpexected error")
		}

		want := string(jsonMarshalled)
		if got := w.Body.String(); got != want {
			t.Errorf("TestJSONWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/json;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestJSONWriter() = %v, want %v", got, wantHeader)
		}
	})

	t.Run("indent", func(t *testing.T) {
		jsonMarshalled, mErr := json.MarshalIndent(sampleObject, "x", "**")
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewJSONIndentWriter("x", "**", emissione.StreamMethod(nil))

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestJSONWriter() unpexected error")
		}

		want := string(jsonMarshalled)
		if got := w.Body.String(); got != want {
			t.Errorf("TestJSONWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/json;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestJSONWriter() = %v, want %v", got, wantHeader)
		}
	})

	t.Run("simple-stream", func(t *testing.T) {
		jsonMarshalled, mErr := json.Marshal(sampleObject)
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewJSONWriter()

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestJSONWriter() unpexected error")
		}

		want := string(jsonMarshalled)
		if got := strings.TrimSpace(w.Body.String()); got != want {
			t.Errorf("TestJSONWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/json;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestJSONWriter() = %v, want %v", got, wantHeader)
		}
	})

	t.Run("indent-stream", func(t *testing.T) {
		jsonMarshalled, mErr := json.MarshalIndent(sampleObject, "x", "**")
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewJSONIndentWriter("x", "**")

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestJSONWriter() unpexected error")
		}

		want := string(jsonMarshalled)
		if got := strings.TrimSpace(w.Body.String()); got != want {
			t.Errorf("TestJSONWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/json;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestJSONWriter() = %v, want %v", got, wantHeader)
		}
	})
}
