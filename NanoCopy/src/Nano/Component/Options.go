package Component

type (
	options struct {
		name     string              // 组件名字
		nameFunc func(string) string // 重命名handler的名字
	}

	// Option 用来自定义handler
	Option func(options *options)
)

// WithName 重命名组件的名字
func WithName(name string) Option {
	return func(opt *options) {
		opt.name = name
	}
}

// WithNameFunc 用指定函数覆盖handler的名字
// 例如: strings.ToUpper/strings.ToLower
func WithNameFunc(fn func(string) string) Option {
	return func(opt *options) {
		opt.nameFunc = fn
	}
}
