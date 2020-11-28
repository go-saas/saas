package data

type SeedContext struct {
	TenantId string
	//extra properties
	P map[string]interface{}
}

func NewSeedContext(tenantId string) *SeedContext {
	return &SeedContext{
		TenantId: tenantId,
		P:        make(map[string]interface{}),
	}
}

func (s *SeedContext) WithKV(k string, v interface{}) *SeedContext {
	s.P[k] = v
	return s
}
