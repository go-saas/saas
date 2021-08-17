package seed

type SeedOption struct {
	Contributors []SeedContributor
	TenantIds    []string
}

func NewSeedOption(opt ...SeedContributor) *SeedOption {
	return &SeedOption{Contributors: opt, TenantIds: make([]string, 0)}
}
