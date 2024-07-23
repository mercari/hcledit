package hcledit

type option struct {
	comment                 string
	afterKey                string
	beforeNewline           bool
	readFallbackToRawString bool
	querySeparator          rune
}

// Option configures specific behavior for specific HCLEditor operations.
// TODO(slewiskelly): Not all options are applicable to all operations, maybe
// options should be specific to each kind of operation?
type Option func(*option)

// WithComment provides comment to put together when creating.
func WithComment(comment string) Option {
	return func(opt *option) {
		opt.comment = comment
	}
}

func WithAfter(key string) Option {
	return func(opt *option) {
		opt.afterKey = key
	}
}

// WithNewLine
func WithNewLine() Option {
	return func(opt *option) {
		opt.beforeNewline = true
	}
}

// This provides a fallback to return the raw string of the value if we could
// not parse it. If this option is provided to HCLEditor.Read(), the error
// return value will signal fallback occurred.
func WithReadFallbackToRawString() Option {
	return func(opt *option) {
		opt.readFallbackToRawString = true
	}
}

// This sets a separator for queryStr in HCLEditor funcs. The default separator is ".".
func WithQuerySeparator(r rune) Option {
	return func(opt *option) {
		opt.querySeparator = r
	}
}
