package saas

type DatabaseStyleType int32

const (
	Single    DatabaseStyleType = 1 << 0
	PerTenant DatabaseStyleType = 1 << 1
	Multi     DatabaseStyleType = 1 << 2
)

type MultiTenancyOption struct {
	IsEnabled     bool
	DatabaseStyle DatabaseStyleType
}

type option func(tenancyOption *MultiTenancyOption)

// WithEnabled enable status
func WithEnabled(isEnabled bool) option {
	return func(tenancyOption *MultiTenancyOption) {
		tenancyOption.IsEnabled = isEnabled
	}
}

//
// WithDatabaseStyle database style, support Single/PerTenant/Multi
func WithDatabaseStyle(databaseStyle DatabaseStyleType) option {
	return func(tenancyOption *MultiTenancyOption) {
		tenancyOption.DatabaseStyle = databaseStyle
	}
}

func NewMultiTenancyOption(opts ...option) *MultiTenancyOption {
	option := MultiTenancyOption{}
	for _, opt := range opts {
		opt(&option)
	}
	return &option
}

func DefaultMultiTenancyOption() *MultiTenancyOption {
	return NewMultiTenancyOption(WithEnabled(true), WithDatabaseStyle(Multi))
}
