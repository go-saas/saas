package http

const DefaultKey = "__tenant"

func KeyOrDefault(key string) string {
	if len(key) > 0 {
		return key

	}
	return DefaultKey
}

type WebMultiTenancyOption struct {
	TenantKey    string
	DomainFormat string
}

func NewWebMultiTenancyOption(key string, domainFormat string) *WebMultiTenancyOption {
	key = KeyOrDefault(key)
	return &WebMultiTenancyOption{
		TenantKey:    key,
		DomainFormat: domainFormat,
	}
}

func NewDefaultWebMultiTenancyOption() *WebMultiTenancyOption {
	return NewWebMultiTenancyOption("", "")
}
