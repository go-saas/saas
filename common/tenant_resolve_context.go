package common

import "context"

type TenantResolveContext struct {
	context.Context
	TenantIdOrName string
	// HasHandled field to handle host side unresolved or resolved
	HasHandled bool
}

func (t TenantResolveContext) HasResolved() bool {
	return t.HasHandled
}

var DefaultTenantResolveContext = TenantResolveContext{
	TenantIdOrName: "",
	HasHandled:     false,
}
