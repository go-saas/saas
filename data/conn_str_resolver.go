package data

import "context"

type ConnStrResolver interface {
	// Resolve connection string by key
	Resolve(ctx context.Context, key string) string
}

// DefaultConnStrResolver use config map to resolve connection string
type DefaultConnStrResolver struct {
	Opt *ConnStrOption
}

func NewDefaultConnStrResolver(c *ConnStrOption) *DefaultConnStrResolver{
	return &DefaultConnStrResolver{
		Opt: c,
	}
}

// Resolve from option
func (d DefaultConnStrResolver) Resolve(_ context.Context, key string) string {
	if key != "" {
		v := d.Opt.Conn[key]
		if v != "" {
			return v
		}
	}
	return d.Opt.Conn.Default()
}
