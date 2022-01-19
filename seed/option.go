package seed

import "github.com/goxiaoy/uow"

type Option struct {
	Contributors []Contributor
	TenantIds    []string
	uowMgr       uow.Manager
}

func NewSeedOption(opt ...Contributor) *Option {
	return &Option{Contributors: opt, TenantIds: make([]string, 0)}
}

func (opt *Option) WithUow(uowMgr uow.Manager) *Option {
	opt.uowMgr = uowMgr
	return opt
}

func (opt *Option) WithTenantId(tenants ...string) *Option {
	opt.TenantIds = tenants
	return opt
}

func (opt *Option) AddTenantId(tenants ...string) *Option {
	opt.TenantIds = append(opt.TenantIds, tenants...)
	return opt
}
