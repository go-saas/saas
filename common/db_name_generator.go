package common

import (
	"context"
	"github.com/goxiaoy/go-saas/data"
)

// DbNameGenerator generate table name for tenant. useful for tenant creation
type DbNameGenerator interface {
	Gen(ctx context.Context, tenant TenantInfo) (string, error)
}

// DefaultDbNameGenerator generate only database name for tenant
type DefaultDbNameGenerator struct {
	prefix  string
	postfix string
}

var _ DbNameGenerator = (*DefaultDbNameGenerator)(nil)

func NewDbNameGenerator(prefix string, postfix string) *DefaultDbNameGenerator {
	return &DefaultDbNameGenerator{prefix: prefix, postfix: postfix}
}

func (d *DefaultDbNameGenerator) Gen(ctx context.Context, tenant TenantInfo) (string, error) {
	return d.prefix + tenant.GetId() + d.postfix, nil
}

func NewSqliteFileDbNameGenerator(parent string) *DefaultDbNameGenerator {
	return NewDbNameGenerator(parent, ".db")
}

// MultiDbNameGenerator generate  data.ConnStrings for tenant
type MultiDbNameGenerator struct {
	//conf host connection string
	conf data.ConnStrings
	g    DbNameGenerator
}

func NewMultipleDbNameGenerator(conf data.ConnStrings, g DbNameGenerator) *MultiDbNameGenerator {
	return &MultiDbNameGenerator{conf: conf, g: g}
}

func (m *MultiDbNameGenerator) GenMulti(ctx context.Context, tenant TenantInfo) (data.ConnStrings, error) {
	ret := data.ConnStrings{}
	name, err := m.g.Gen(ctx, tenant)
	if err != nil {
		return nil, err
	}
	if m.conf != nil {
		for k, _ := range m.conf {
			ret[k] = name
		}
	}
	ret[data.Default] = name
	return ret, nil
}
