package emissione

// WriterOption is the option-wrapper for defining the workings of an Writer
type WriterOptions struct {
	MarshallMethod func(v interface{}) ([]byte, error)
	ContentType    string
}

type WriterOption func(*WriterOptions)

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
