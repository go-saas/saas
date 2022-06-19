package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	mysql2 "github.com/go-sql-driver/mysql"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gin/saas"
	gorm2 "github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/seed"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
)

const (
	defaultSqliteSharedDsn = "./shared.db"
	defaultMysqlSharedDsn  = "root:youShouldChangeThis@tcp(127.0.0.1:3306)/shared?parseTime=true&loc=Local"
)

var (
	driver        string
	sharedDsn     string
	tenant3Dsn    string
	ensureDbExist func(s string) error
)

func init() {
	flag.StringVar(&driver, "driver", "sqlite3", "sqlite3/mysql")
	flag.StringVar(&sharedDsn, "dsn", "", "shared dsn.")
}

func main() {
	flag.Parse()

	switch driver {
	case sqlite.DriverName:
		if len(sharedDsn) == 0 {
			sharedDsn = defaultSqliteSharedDsn
			tenant3Dsn = "./tenant3.db"
		}
	case "mysql":
		if len(sharedDsn) == 0 {
			sharedDsn = defaultMysqlSharedDsn
			tenant3, err := mysql2.ParseDSN(sharedDsn)
			if err != nil {
				panic(err)
			}
			tenant3.DBName = "tenant3"
			tenant3Dsn = tenant3.FormatDSN()

			ensureDbExist = func(s string) error {
				dsn, err := mysql2.ParseDSN(s)
				if err != nil {
					return err
				}
				dbname := dsn.DBName
				dsn.DBName = ""
				//open without db name
				db, err := g.Open(mysql.Open(dsn.FormatDSN()))
				if err != nil {
					return err
				}
				err = db.Debug().Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbname)).Error
				if err != nil {
					return err
				}
				return closeDb(db)
			}

		}
	default:
		panic(fmt.Errorf("driver %s unsupported", driver))
	}

	r := gin.Default()
	//dbOpener
	dbOpener, c := common.NewCachedDbOpener(common.DbOpenerFunc(func(s string) (*sql.DB, error) {
		if ensureDbExist != nil {
			if err := ensureDbExist(s); err != nil {
				return nil, err
			}
		}
		return sql.Open(driver, s)
	}))
	defer c()
	clientProvider := gorm2.ClientProviderFunc(func(ctx context.Context, s string) (*g.DB, error) {
		db, err := dbOpener.Open(s)
		if err != nil {
			return nil, err
		}
		var client *g.DB

		if driver == sqlite.DriverName {
			db.SetMaxIdleConns(1)
			db.SetMaxOpenConns(1)

			client, err = g.Open(&sqlite.Dialector{
				DriverName: sqlite.DriverName,
				DSN:        s,
				Conn:       db,
			})
			if err != nil {
				return client, err
			}
		} else if driver == "mysql" {
			client, err = g.Open(mysql.New(mysql.Config{
				Conn: db,
			}))
			if err != nil {
				return client, err
			}
		}

		return client.WithContext(ctx).Debug(), err
	})

	conn := make(data.ConnStrings, 1)
	//default database
	conn.SetDefault(sharedDsn)

	var tenantStore common.TenantStore

	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return tenantStore
	}, data.NewConnStrOption(conn))
	dbProvider := gorm2.NewDbProvider(mr, clientProvider)

	tenantStore = common.NewCachedTenantStore(&TenantStore{dbProvider: dbProvider})

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

func closeDb(d *g.DB) error {
	sqlDB, err := d.DB()
	if err != nil {
		return err
	}
	cErr := sqlDB.Close()
	if cErr != nil {
		//todo logging
		//logger.Errorf("Gorm db close error: %s", err.Error())
		return cErr
	}
	return nil
}
