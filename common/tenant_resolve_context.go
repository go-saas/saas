package common

type TenantResolveContext struct {
	TenantIdOrName string
	// HasHandled field to handle host side unresolved or resolved
	HasHandled bool
}

func (t TenantResolveContext) HasResolved() bool {
	return t.TenantIdOrName != "" || t.HasHandled
}

var DefaultTenantResolveContext = TenantResolveContext{
	TenantIdOrName: "",
	HasHandled:     false,
}
