// Code generated by entc, DO NOT EDIT.

package runtime

import (
	"context"

	"github.com/goxiaoy/go-saas/examples/ent/tenant/ent/post"
	"github.com/goxiaoy/go-saas/examples/ent/tenant/ent/schema"

	"entgo.io/ent"
	"entgo.io/ent/privacy"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	postMixin := schema.Post{}.Mixin()
	post.Policy = privacy.NewPolicies(postMixin[0], schema.Post{})
	post.Hooks[0] = func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if err := post.Policy.EvalMutation(ctx, m); err != nil {
				return nil, err
			}
			return next.Mutate(ctx, m)
		})
	}
}

const (
	Version = "v0.10.1"                                         // Version of ent codegen.
	Sum     = "h1:dM5h4Zk6yHGIgw4dCqVzGw3nWgpGYJiV4/kyHEF6PFo=" // Sum of ent codegen.
)
