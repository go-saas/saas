package seed

import "context"

type Contributor interface {
	Seed(ctx context.Context, sCtx *Context) error
}
