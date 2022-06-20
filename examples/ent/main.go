package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	sent "github.com/goxiaoy/go-saas/ent"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent"
	_ "github.com/goxiaoy/go-saas/examples/ent/shared/ent/runtime"
	ent2 "github.com/goxiaoy/go-saas/examples/ent/tenant/ent"
	_ "github.com/goxiaoy/go-saas/examples/ent/tenant/ent/runtime"
	"github.com/goxiaoy/go-saas/gin/saas"
	"github.com/goxiaoy/go-saas/seed"
	_ "github.com/mattn/go-sqlite3"
)

type SharedDbProvider common.DbProvider[*ent.Client]
type TenantDbProvider common.DbProvider[*ent2.Client]

func main() {
	r := gin.Default()

	cache := common.NewCache[string, *ent.Client]()
	defer cache.Flush()
	cache2 := common.NewCache[string, *ent2.Client]()
	defer cache.Flush()

	sharedClientProvider := common.ClientProviderFunc[*ent.Client](func(ctx context.Context, s string) (*ent.Client, error) {
		v, _, err := cache.GetOrSet(s, func() (*ent.Client, error) {
			client, err := ent.Open("sqlite3", s, ent.Debug())
			if err != nil {
				return nil, err
			}
			client.Use(sent.Saas)
			return client, nil
		})
		return v, err
	})
	tenantClientProvider := common.ClientProviderFunc[*ent2.Client](func(ctx context.Context, s string) (*ent2.Client, error) {

		v, _, err := cache2.GetOrSet(s, func() (*ent2.Client, error) {
			client, err := ent2.Open("sqlite3", s, ent2.Debug())
			if err != nil {
				return nil, err
			}
			client.Use(sent.Saas)
			return client, nil
		})
		return v, err
	})

	conn := make(data.ConnStrings, 1)
	//default database
	conn.SetDefault("./shared.db?_fk=1")

	var tenantStore common.TenantStore

	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return tenantStore
	}, data.NewConnStrOption(conn))

	sharedDbProvider := common.NewDbProvider[*ent.Client](mr, sharedClientProvider)
	tenantDbProvider := common.NewDbProvider[*ent2.Client](mr, tenantClientProvider)

	tenantStore = common.NewCachedTenantStore(&TenantStore{shared: sharedDbProvider})

	r.Use(saas.MultiTenancy(tenantStore))

	//return current tenant
	r.GET("/tenant/current", func(c *gin.Context) {
		rCtx := c.Request.Context()
		tenantInfo, _ := common.FromCurrentTenant(rCtx)
		trR := common.FromTenantResolveRes(rCtx)
		c.JSON(200, gin.H{
			"tenantId":  tenantInfo.GetId(),
			"resolvers": trR.AppliedResolvers,
		})
	})

	r.GET("/posts", func(c *gin.Context) {
		ctx := c.Request.Context()
		tenantInfo, _ := common.FromCurrentTenant(ctx)
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
	seedOpt := seed.NewSeedOption().WithTenantId("", "1", "2", "3").WithExtra(map[string]interface{}{})
	err := seeder.Seed(context.Background(), seedOpt)
	if err != nil {
		panic(err)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
