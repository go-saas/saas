package saas

import "context"

type Context struct {
	context        context.Context
	TenantIdOrName string
	// HasHandled field to handle host side unresolved or resolved
	HasHandled bool
}

func NewTenantResolveContext(ctx context.Context) *Context {
	return &Context{
		context: ctx,
	}
}

func (t *Context) HasResolved() bool {
	return t.HasHandled
}

func (t *Context) Context() context.Context {
	return t.context
}

func (t *Context) WithContext(ctx context.Context) {
	t.context = ctx
}
