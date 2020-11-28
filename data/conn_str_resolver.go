package data

import "context"

type ConnStrResolver interface {
	//resolve connection string by key
	Resolve(ctx context.Context, key string) string
}
