package saas

import "github.com/go-saas/saas/data"

type TenantConfig struct {
	ID      string           `json:"id"`
	Name    string           `json:"name"`
	Region  string           `json:"region"`
	PlanKey string           `json:"planKey"`
	Conn    data.ConnStrings `json:"conn"`
}

func NewTenantConfig(id, name, region, planKey string) *TenantConfig {
	return &TenantConfig{
		ID:      id,
		Name:    name,
		Region:  region,
		PlanKey: planKey,
		Conn:    make(data.ConnStrings),
	}
}
