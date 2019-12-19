package emissione_test

import (
	"github.com/kernle32dll/emissione-go"

	"encoding/xml"
	"net/http/httptest"
	"testing"
)

func TestXMLWriter(t *testing.T) {
	type SampleObject struct {
		TestProp int `xml:"TestProp"`
	}

	sampleObject := SampleObject{42}

	t.Run("simple", func(t *testing.T) {
		xmlMarshalled, mErr := xml.Marshal(sampleObject)
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewXmlWriter(emissione.StreamMethod(nil))

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestXMLWriter() unpexected error")
		}

		want := string(xmlMarshalled)
		if got := w.Body.String(); got != want {
			t.Errorf("TestXMLWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/xml;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestXMLWriter() = %v, want %v", got, wantHeader)
		}
	})

	t.Run("indent", func(t *testing.T) {
		xmlMarshalled, mErr := xml.MarshalIndent(sampleObject, "x", "**")
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewXmlIndentWriter("x", "**", emissione.StreamMethod(nil))

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestXMLWriter() unpexected error")
		}

		want := string(xmlMarshalled)
		if got := w.Body.String(); got != want {
			t.Errorf("TestXMLWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/xml;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestXMLWriter() = %v, want %v", got, wantHeader)
		}
	})

	t.Run("simple-stream", func(t *testing.T) {
		xmlMarshalled, mErr := xml.Marshal(sampleObject)
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewXmlWriter()

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestXMLWriter() unpexected error")
		}

		want := string(xmlMarshalled)
		if got := w.Body.String(); got != want {
			t.Errorf("TestXMLWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/xml;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestXMLWriter() = %v, want %v", got, wantHeader)
		}
	})

	t.Run("indent-stream", func(t *testing.T) {
		xmlMarshalled, mErr := xml.MarshalIndent(sampleObject, "x", "**")
		if mErr != nil {
			t.Error(mErr)
		}

		writer := emissione.NewXmlIndentWriter("x", "**")

		w := httptest.NewRecorder()

		if writer.Write(w, sampleObject) != nil {
			t.Errorf("TestXMLWriter() unpexected error")
		}

		want := string(xmlMarshalled)
		if got := w.Body.String(); got != want {
			t.Errorf("TestXMLWriter() = %v, want %v", got, want)
		}

		wantHeader := "application/xml;charset=utf-8"
		if got := w.Header().Get("Content-Type"); got != wantHeader {
			t.Errorf("TestXMLWriter() = %v, want %v", got, wantHeader)
		}
	})
}
