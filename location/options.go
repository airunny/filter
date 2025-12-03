package location

type Options struct {
	Language string
}

type Option func(*Options)

func WithLanguage(language string) Option {
	return func(o *Options) {
		o.Language = language
	}
}
