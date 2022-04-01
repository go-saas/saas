package seed

type Option struct {
	TenantIds []string
	Extra     map[string]interface{}
}

func NewSeedOption() *Option {
	return &Option{}
}

func (opt *Option) WithTenantId(tenants ...string) *Option {
	opt.TenantIds = tenants
	return opt
}

func (opt *Option) AddTenantId(tenants ...string) *Option {
	opt.TenantIds = append(opt.TenantIds, tenants...)
	return opt
}

func (opt *Option) WithExtra(extra map[string]interface{}) *Option {
	opt.Extra = extra
	return opt
}
