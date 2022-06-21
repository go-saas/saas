package data

import "context"

type ConnStrResolver interface {
	// Resolve connection string by user-friendly key
	Resolve(ctx context.Context, key string) (string, error)
}

type ConnStrResolverFunc func(ctx context.Context, key string) (string, error)

func (c ConnStrResolverFunc) Resolve(ctx context.Context, key string) (string, error) {
	return c(ctx, key)
}

var _ ConnStrResolver = (*ConnStrResolverFunc)(nil)

func ChainConnStrResolver(cs ...ConnStrResolver) ConnStrResolver {
	return ConnStrResolverFunc(func(ctx context.Context, key string) (string, error) {
		for _, c := range cs {
			conn, err := c.Resolve(ctx, key)
			if err != nil {
				return "", err
			}
			if len(conn) > 0 {
				return conn, err
			}
		}
		return "", nil
	})
}
