package emissione

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var (
	defHandler = Default()
)

type Writer interface {
	Write(w http.ResponseWriter, i interface{}) error
}

type Handler struct {
	defaultHandler Writer
	handlers       map[string]Writer

	wildcardDetector *regexp.Regexp
}

func Default() *Handler {
	json := NewJSONIndentWriter("", "  ")
	xml := NewXmlIndentWriter("", "  ")

	return New(
		json,
		map[string]Writer{
			"application/json":                json,
			"application/json;charset=utf-8":  json,
			"application/json; charset=utf-8": json,
			"application/xml":                 xml,
			"application/xml;charset=utf-8":   xml,
			"application/xml; charset=utf-8":  xml,
		},
	)
}

func New(defaultHandler Writer, handlerMapping map[string]Writer) *Handler {
	wildcardDetector, err := regexp.Compile(".*?/\\*(;q=[\\d.]+)?")
	if err != nil {
		panic(fmt.Sprintf("unexpected error building regex: %s", err))
	}

	lowerCaseMapping := make(map[string]Writer, len(handlerMapping))
	for k, v := range handlerMapping {
		lowerCaseMapping[strings.ToLower(k)] = v
	}

	return &Handler{
		wildcardDetector: wildcardDetector,

		defaultHandler: defaultHandler,
		handlers:       lowerCaseMapping,
	}
}

func (h Handler) findWriterByRequest(r *http.Request) Writer {
	acceptHeader := r.Header.Get("Accept")

	// If Accept headers was not set, use the default
	if acceptHeader == "" {
		return h.defaultHandler
	}

	// Split and lower-case
	accepts := strings.Split(acceptHeader, ",")
	for i := range accepts {
		accepts[i] = strings.ToLower(accepts[i])
	}

	acceptsPrioritized := AcceptSlice(accepts)
	sort.Sort(sort.Reverse(acceptsPrioritized))

	// Iterate all send types
	for _, accept := range acceptsPrioritized {
		if possibleWriter := h.findWriterByType(accept); possibleWriter != nil {
			return possibleWriter
		}
	}

	return h.defaultHandler
}

func (h Handler) findWriterByType(acceptType string) Writer {
	// 1: Wildcard?
	if acceptType == "*/*" {
		return h.defaultHandler
	}

	// 2: Exact match?
	if writer, match := h.handlers[acceptType]; match {
		return writer
	}

	// 3: Exact match without quality?
	if strings.Contains(acceptType, "q=") {
		accepts := strings.Split(acceptType, "q=")

		acceptQualityLess := accepts[0]
		if len(accepts) > 1 {
			acceptQualityLess = strings.Join(accepts[:len(accepts)-1], "")
		}

		acceptQualityLess = strings.Trim(acceptQualityLess, "; ")

		if writer, match := h.handlers[acceptQualityLess]; match {
			return writer
		}
	}

	// 4: Partial wildcard? (e.g. image/*)
	if h.wildcardDetector.MatchString(acceptType) {
		return h.findPartialWriterMatch(strings.Split(acceptType, "/")[0])
	}

	return nil
}

func (h Handler) findPartialWriterMatch(mimeType string) Writer {
	applicableHandlers := map[string]Writer{}

	// Find all handlers which are applicable
	for k, v := range h.handlers {
		if strings.HasPrefix(k, mimeType) {
			applicableHandlers[k] = v
		}
	}

	handlerKeys := make([]string, len(applicableHandlers))
	sort.Strings(handlerKeys)

	return applicableHandlers[handlerKeys[0]]
}

func (h Handler) Write(w http.ResponseWriter, r *http.Request, code int, i interface{}) {
	w.WriteHeader(code)

	if err := h.findWriterByRequest(r).Write(w, i); err != nil {
		panic(err)
	}
}

func Write(w http.ResponseWriter, r *http.Request, code int, i interface{}) {
	defHandler.Write(w, r, code, i)
}
