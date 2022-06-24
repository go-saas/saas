package seed

type SeedOption struct {
	TenantIds []string
	Extra     map[string]interface{}
}

func NewOption() *SeedOption {
	return &SeedOption{Extra: map[string]interface{}{}}
}

type Option func(opt *SeedOption)

func WithTenantId(tenants ...string) Option {
	return func(opt *SeedOption) {
		opt.TenantIds = tenants
	}
}

func AddHost() Option {
	return func(opt *SeedOption) {
		opt.TenantIds = append(opt.TenantIds, "")
	}
}
func AddTenant(tenants ...string) Option {
	return func(opt *SeedOption) {
		opt.TenantIds = append(opt.TenantIds, tenants...)
	}
}

func WithExtra(extra map[string]interface{}) Option {
	return func(opt *SeedOption) {
		opt.Extra = extra
	}
}

func SetExtra(key string, v interface{}) Option {
	return func(opt *SeedOption) {
		opt.Extra[key] = v
	}
}
