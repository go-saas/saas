package main

import (
	"context"
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
	dbOpener, c := gorm2.NewDbOpener()
	defer c()
	i := 1
	cfg := &gorm2.Config{
		Debug: true,
		Dialect: func(s string) g.Dialector {
			// if use mysql
			// return mysql.Open(s)
			return sqlite.Open(s)
		},
		Cfg: &g.Config{},
		//https://github.com/go-gorm/gorm/issues/2875
		MaxOpenConn: &i,
		MaxIdleConn: &i,
	}

	//create a unit of work manager. you can skip this if you do not want it
	//um := uow.NewManager(&uow.Config{SupportNestedTransaction: false}, func(ctx context.Context, kind, key string) uow.TransactionalDb {
	//	if kind == gorm2.DbKind {
	//		db, err := dbOpener.Open(cfg, key)
	//		if err != nil {
	//			panic(err)
	//		}
	//		return gorm.NewTransactionDb(db)
	//	}
	//	panic(errors.New(fmt.Sprintf("can not resolve %s", key)))
	//})
	//
	//r.Use(Uow(um))

	wOpt := shttp.NewDefaultWebMultiTenancyOption()

	conn := make(data.ConnStrings, 1)
	//default database
	conn.SetDefault("./shared.db")

	tenantStore := common.NewMemoryTenantStore(
		[]common.TenantConfig{
			{ID: "1", Name: "Test1"}, // use default shared.db
			{ID: "2", Name: "Test2"},
			{ID: "3", Name: "Test3", Conn: map[string]string{
				data.Default: "./tenant3.db", // use tenant3.db
			}}},
	)

	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return tenantStore
	}, data.NewConnStrOption(conn))
	dbProvider := gorm2.NewDefaultDbProvider(mr, cfg, dbOpener)

	r.Use(saas.MultiTenancy(wOpt, tenantStore))

	//return current tenant
	r.GET("/tenant/current", func(c *gin.Context) {
		rCtx := c.Request.Context()
		tenantInfo := common.FromCurrentTenant(rCtx)
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
	seeder := seed.NewDefaultSeeder(seed.NewSeedOption(
		NewMigrationSeeder(dbProvider),
		NewSeed(dbProvider)).
		WithTenantId("", "1", "2", "3"),
		//WithUow(um),
		map[string]interface{}{})
	err := seeder.Seed(context.Background())
	if err != nil {
		panic(err)
	}

	r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
