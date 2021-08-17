package seed

type Context struct {
	TenantId string
	//extra properties
	Extra map[string]interface{}
}

func NewSeedContext(tenantId string, extra map[string]interface{}) *Context {
	return &Context{
		TenantId: tenantId,
		Extra:    extra,
	}
}

func (s *Context) WithExtra(k string, v interface{}) *Context {
	s.Extra[k] = v
	return s
}
