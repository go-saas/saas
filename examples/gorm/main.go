package main

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gin/saas"
	gorm2 "github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/seed"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
)

func main() {
	r := gin.Default()
	//dbOpener
	dbOpener, c := common.NewCachedDbOpener(common.DbOpenerFunc(func(s string) (*sql.DB, error) {
		return sql.Open("sqlite3", s)
	}))
	defer c()
	clientProvider := gorm2.ClientProviderFunc(func(ctx context.Context, s string) (*g.DB, error) {
		db, err := dbOpener.Open(s)
		if err != nil {
			return nil, err
		}
		db.SetMaxIdleConns(1)
		db.SetMaxOpenConns(1)

		client, err := g.Open(&sqlite.Dialector{
			DriverName: sqlite.DriverName,
			DSN:        s,
			Conn:       db,
		})
		if err != nil {
			return client, err
		}
		return client.WithContext(ctx), err
	})

	conn := make(data.ConnStrings, 1)
	//default database
	conn.SetDefault("./shared.db")

	var tenantStore common.TenantStore

	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return tenantStore
	}, data.NewConnStrOption(conn))
	dbProvider := gorm2.NewDbProvider(mr, clientProvider)

	tenantStore = common.NewCachedTenantStore(&TenantStore{dbProvider: dbProvider})

	wOpt := shttp.NewDefaultWebMultiTenancyOption()
	r.Use(saas.MultiTenancy(wOpt, tenantStore, nil))

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
		db := dbProvider.Get(c.Request.Context(), "")
		var entities []Post
		if err := db.Model(&Post{}).Find(&entities).Error; err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, entities)
		}
	})

	//seed data into db
	seeder := seed.NewDefaultSeeder(NewMigrationSeeder(dbProvider), NewSeed(dbProvider))
	seedOpt := seed.NewSeedOption().WithTenantId("", "1", "2", "3").WithExtra(map[string]interface{}{})
	err := seeder.Seed(context.Background(), seedOpt)
	if err != nil {
		panic(err)
	}

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
