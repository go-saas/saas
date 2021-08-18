package seed

type Option struct {
	Contributors []Contributor
	TenantIds    []string
}

func NewSeedOption(opt ...Contributor) *Option {
	return &Option{Contributors: opt, TenantIds: make([]string, 0)}
}
