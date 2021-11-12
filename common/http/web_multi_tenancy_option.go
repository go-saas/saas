package http

type WebMultiTenancyOption struct {
	TenantKey    string
	DomainFormat string
}

func NewWebMultiTenancyOption(key string, domainFormat string) *WebMultiTenancyOption {
	if key == "" {
		key = "__tenant"
	}
	return &WebMultiTenancyOption{
		TenantKey:    key,
		DomainFormat: domainFormat,
	}
}

func NewDefaultWebMultiTenancyOption() *WebMultiTenancyOption {
	return NewWebMultiTenancyOption("", "")
}
