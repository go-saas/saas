# [Ent](https://entgo.io/) adapter

- Enable [EntQL Filtering](https://entgo.io/docs/feature-flags/#entql-filtering) and [Privacy Layer](https://entgo.io/docs/feature-flags/#privacy-layer) features
  Modify your `ent/generate.go`
  ```
  go generate ... --feature privacy --feature entql ...
  ```
- Add global runtime hook to your client
  ```go
  import sent "github.com/go-saas/saas/ent"
  client.use(sent.Saas)
  ```
- Copy mixin into your schema
  ```go
  type HasTenant struct {
      mixin.Schema
  }
  
  func (HasTenant) Fields() []ent.Field {
      return []ent.Field{
          field.String("tenant_id").Optional().GoType(&sql.NullString{}),
      }
  }
  
  func (HasTenant) Policy() ent.Policy {
      return privacy.Policy{
          Query: privacy.QueryPolicy{
              FilterTenantRule(),
          },
      }
  }
  
  func FilterTenantRule() privacy.QueryMutationRule {
      type hasTenant interface {
          Where(p entql.P)
          WhereTenantID(p entql.StringP)
      }
      return privacy.FilterFunc(func(ctx context.Context, f privacy.Filter) error {
          ct, _ := common.FromCurrentTenant(ctx)
          e := data.FromMultiTenancyDataFilter(ctx)
          hf, ok := f.(hasTenant)
          if e && ok {
              //apply data filter
              if ct.GetId() == "" {
                  //host side
                  hf.Where(entql.FieldNil("tenant_id"))
              } else {
                  hf.WhereTenantID(entql.StringEQ(ct.GetId()))
              }
          }
  
          return privacy.Skip
      })
  }
    
  ```
- Embed mixin into your schema
  ```go
  // Post holds the schema definition for the Post entity.
  type Post struct {
      ent.Schema
  }
  
  func (Post) Mixin() []ent.Mixin {
      return []ent.Mixin{
          HasTenant{},
      }
  }
  ```