package hcledit

type option struct {
	comment       string
	afterKey      string
	beforeNewline bool
}

type Option func(*option)

// WithComment provides comment to put together when creating.
func WithComment(comment string) Option {
	return func(opt *option) {
		opt.comment = comment
	}
}

// WithAfter
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
