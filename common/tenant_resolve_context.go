package common

import "context"

type TenantResolveContext struct {
	context        context.Context
	TenantIdOrName string
	// HasHandled field to handle host side unresolved or resolved
	HasHandled bool
}

func NewTenantResolveContext(ctx context.Context) *TenantResolveContext {
	return &TenantResolveContext{
		context: ctx,
	}
}

func (t *TenantResolveContext) HasResolved() bool {
	return t.HasHandled
}

func (t *TenantResolveContext) Context() context.Context {
	return t.context
}

func (t *TenantResolveContext) WithContext(ctx context.Context) *TenantResolveContext {
	t.context = ctx
	return t
}
