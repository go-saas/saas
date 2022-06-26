package apisix

import (
	"context"
	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/go-saas/saas"

	"regexp"
)

type Resolver struct {
	r         pkgHTTP.Request
	key       string
	pathRegex string
}

func NewResolver(r pkgHTTP.Request, key, pathRegex string) *Resolver {
	return &Resolver{
		r:         r,
		key:       key,
		pathRegex: pathRegex,
	}
}

var _ saas.TenantResolver = (*Resolver)(nil)

func (r *Resolver) Resolve(ctx context.Context) (saas.TenantResolveResult, context.Context, error) {
	// default host side
	var t = ""
	if v := r.r.Header().Get(r.key); len(v) > 0 {
		t = v
	}
	if v := r.r.Args().Get(r.key); len(v) > 0 {
		t = v
	}
	if len(r.pathRegex) > 0 {
		reg := regexp.MustCompile(r.pathRegex)
		f := reg.FindAllStringSubmatch(string(r.r.Path()), -1)
		if f != nil {
			t = f[0][1]
		}
	}

	return saas.TenantResolveResult{TenantIdOrName: t}, ctx, nil
}
