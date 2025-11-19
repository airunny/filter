package variable

type options struct {
	name string
}

type Option func(o *options)

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

type SimpleBuilder struct {
	name     string
	variable Variable
}

func (s *SimpleBuilder) Name() string {
	if s.name != "" {
		return s.name
	}
	return s.variable.Name()
}

func (s *SimpleBuilder) Build(_ string) Variable {
	return s.variable
}

func NewSimpleVariable(variable Variable, opts ...Option) *SimpleBuilder {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return &SimpleBuilder{
		name:     o.name,
		variable: variable,
	}
}
