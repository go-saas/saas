# [Ent](https://entgo.io/) adapter

- Enable [EntQL Filtering](https://entgo.io/docs/feature-flags/#entql-filtering) and [Privacy Layer](https://entgo.io/docs/feature-flags/#privacy-layer) features
  Modify your `ent/generate.go`
  ```
  go generate ... --feature intercept,schema/snapshot ...
  ```


- Embed mixin into your schema

  ```go
  import (
  	sent "github.com/go-saas/saas/ent"
  )
  ...
  // Post holds the schema definition for the Post entity.
  type Post struct {
      ent.Schema
  }
  
  func (Post) Mixin() []ent.Mixin {
      return []ent.Mixin{
          sent.HasTenant{},
      }
  }
  ```
