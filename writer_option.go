package emissione

// WriterOptions is the option-wrapper for defining the workings of a Writer.
type WriterOptions struct {
	MarshallMethod func(v interface{}) ([]byte, error)
	ContentType    string
}

// WriterOption is functional option to WriterOptions.
type WriterOption func(*WriterOptions)

// MarshallMethod defines the marshalling method used by a SimpleWriter
// to marshall output objects.
func MarshallMethod(method func(v interface{}) ([]byte, error)) WriterOption {
	return func(args *WriterOptions) {
		args.MarshallMethod = method
	}
}

// ContentType sets the content-type string, which will be set as the header
// of the same name.
func ContentType(contentType string) WriterOption {
	return func(args *WriterOptions) {
		args.ContentType = contentType
	}
}
