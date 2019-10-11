package emissione

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Write(t *testing.T) {
	type SampleObject struct {
		TestProp int `json:"testProp",xml:"testProp"`
	}

	sampleObject := SampleObject{42}

	jsonIndent, err := json.MarshalIndent(sampleObject, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	xmlIndent, err := xml.MarshalIndent(sampleObject, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		accepts []string
		want    string
	}{
		{"default", []string{"unknown"}, string(jsonIndent)},
		{"empty-header", nil, string(jsonIndent)},
		{"empty-but-different", []string{""}, string(jsonIndent)},
		{"*/*", []string{"*/*"}, string(jsonIndent)},
		{"application/json", []string{"application/json"}, string(jsonIndent)},
		{"application/*", []string{"application/*"}, string(jsonIndent)},
		{"application/xml", []string{"application/xml"}, string(xmlIndent)},
		{"weighted", []string{"application/*, application/xml;q=2"}, string(xmlIndent)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			Write(w, &http.Request{Header: map[string][]string{"Accept": tt.accepts}}, http.StatusTeapot, sampleObject)

			if got := w.Body.String(); got != tt.want {
				t.Errorf("Write() = %v, want %v", got, tt.want)
			}
		})
	}
}
