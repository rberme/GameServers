package Nano

type (
	options struct {
		pipeline IPipeline
	}

	Option func(*options)
)

func WithPipeline(pipeline IPipeline) Option {
	return func(opt *options) {
		opt.pipeline = pipeline
	}
}
