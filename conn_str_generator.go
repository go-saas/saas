package saas

import (
	"context"
	"fmt"
)

// ConnStrGenerator generate connection string for tenant. useful for tenant creation
type ConnStrGenerator interface {
	Gen(ctx context.Context, tenant TenantInfo) (string, error)
}

type DefaultConnStrGenerator struct {
	format string
}

var _ ConnStrGenerator = (*DefaultConnStrGenerator)(nil)

func NewConnStrGenerator(format string) *DefaultConnStrGenerator {
	return &DefaultConnStrGenerator{format: format}
}

func (d *DefaultConnStrGenerator) Gen(ctx context.Context, tenant TenantInfo) (string, error) {
	if len(tenant.GetId()) == 0 {
		return "", nil
	}
	return fmt.Sprintf(d.format, tenant.GetId()), nil
}
