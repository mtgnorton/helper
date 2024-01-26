package statistics

type Options struct {
	saveResult bool
}
type Option func(o *Options)

func WithSaveResult(saveResult bool) Option {
	return func(o *Options) {
		o.saveResult = saveResult
	}
}
