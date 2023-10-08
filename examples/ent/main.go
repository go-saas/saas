package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-saas/saas"
	"github.com/go-saas/saas/data"
	"github.com/go-saas/saas/examples/ent/shared/ent"
	_ "github.com/go-saas/saas/examples/ent/shared/ent/runtime"
	ent2 "github.com/go-saas/saas/examples/ent/tenant/ent"
	_ "github.com/go-saas/saas/examples/ent/tenant/ent/runtime"
	sgin "github.com/go-saas/saas/gin"
	"github.com/go-saas/saas/seed"
	_ "github.com/mattn/go-sqlite3"
)

type SharedDbProvider saas.DbProvider[*ent.Client]
type TenantDbProvider saas.DbProvider[*ent2.Client]

func main() {
	r := gin.Default()

	cache := saas.NewCache[string, *ent.Client]()
	defer cache.Flush()
	cache2 := saas.NewCache[string, *ent2.Client]()
	defer cache.Flush()

	sharedClientProvider := saas.ClientProviderFunc[*ent.Client](func(ctx context.Context, s string) (*ent.Client, error) {
		v, _, err := cache.GetOrSet(s, func() (*ent.Client, error) {
			client, err := ent.Open("sqlite3", s, ent.Debug())
			if err != nil {
				return nil, err
			}
			return client, nil
		})
		return v, err
	})
	tenantClientProvider := saas.ClientProviderFunc[*ent2.Client](func(ctx context.Context, s string) (*ent2.Client, error) {
		v, _, err := cache2.GetOrSet(s, func() (*ent2.Client, error) {
			client, err := ent2.Open("sqlite3", s, ent2.Debug())
			if err != nil {
				return nil, err
			}
			return client, nil
		})
		return v, err
	})

	conn := make(data.ConnStrings, 1)
	//default database
	conn.SetDefault("./shared.db?_fk=1")

	var tenantStore saas.TenantStore

	//host (shared) database use connection string from config
	sharedDbProvider := saas.NewDbProvider[*ent.Client](conn, sharedClientProvider)

	tenantStore = &TenantStore{shared: sharedDbProvider}

	mr := saas.NewMultiTenancyConnStrResolver(tenantStore, conn)
	// tenant database use connection string from tenantStore
	tenantDbProvider := saas.NewDbProvider[*ent2.Client](mr, tenantClientProvider)

	r.Use(sgin.MultiTenancy(tenantStore))

	//return current tenant
	r.GET("/tenant/current", func(c *gin.Context) {
		rCtx := c.Request.Context()
		tenantInfo, _ := saas.FromCurrentTenant(rCtx)
		trR := saas.FromTenantResolveRes(rCtx)
		c.JSON(200, gin.H{
			"tenantId":  tenantInfo.GetId(),
			"resolvers": trR.AppliedResolvers,
		})
	})

	r.GET("/posts", func(c *gin.Context) {
		ctx := c.Request.Context()
		tenantInfo, _ := saas.FromCurrentTenant(ctx)
		if tenantInfo.GetId() == "" {
			db := sharedDbProvider.Get(ctx, "")
			e, err := db.Post.Query().All(ctx)
			if err != nil {
				c.AbortWithError(500, err)
			}
			c.JSON(200, e)
		} else {
			db := tenantDbProvider.Get(ctx, "")
			e, err := db.Post.Query().All(ctx)
			if err != nil {
				c.AbortWithError(500, err)
			}
			c.JSON(200, e)
		}
	})

	//seed data into db
	seeder := seed.NewDefaultSeeder(NewMigrationSeeder(sharedDbProvider, tenantDbProvider), NewSeed(sharedDbProvider, tenantDbProvider))
	err := seeder.Seed(context.Background(), seed.AddHost(), seed.AddTenant("1", "2", "3"))
	if err != nil {
		panic(err)
	}

	r.Run(":8090") // listen and serve on 0.0.0.0:8090 (for windows "localhost:8090")
}
