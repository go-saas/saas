package data

type SeedOption struct {
	Contributors []SeedContributor
}

func NewSeedOption(opt ...SeedContributor) *SeedOption {
	return &SeedOption{Contributors: opt}
}
