package http

type PatchHttpMultiTenancyOption func(tenancyOption *WebMultiTenancyOption)

type WebMultiTenancyOption struct {
	TenantKey string
}

func NewWebMultiTenancyOption(key string) *WebMultiTenancyOption {
	if key ==""{
		key = "__tenant"
	}
	return &WebMultiTenancyOption{
		TenantKey: key,
	}
}
func DefaultWebMultiTenancyOption() *WebMultiTenancyOption {
	return NewWebMultiTenancyOption("")
}
