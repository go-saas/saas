package http

type PatchHttpMultiTenancyOption func(tenancyOption *MultiTenancyOption)

type MultiTenancyOption struct {
	TenantKey string
}

func NewMultiTenancyOption(key string) *MultiTenancyOption {
	if key ==""{
		key = "__tenant"
	}
	return &MultiTenancyOption{
		TenantKey: key,
	}
}
