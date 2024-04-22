package hcledit

type option struct {
	comment                 string
	afterKey                string
	beforeNewline           bool
	readFallbackToRawString bool
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

func WithReadFallbackToRawString() Option {
	return func(opt *option) {
		opt.readFallbackToRawString = true
	}
}
