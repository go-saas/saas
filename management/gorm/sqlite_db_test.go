package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	gorm2 "github.com/goxiaoy/uow/gorm"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"testing"
)

var TestDb *g.DB
var TestDbProvider *gorm.DefaultDbProvider
var TestTenantRepo *TenantRepo
var TestGormTenantStore *TenantStore
var c gorm.DbClean
var TestDbOpener gorm.DbOpener

var TestUnitOfWorkManager uow.Manager

func TestMain(m *testing.M) {
	TestDbOpener, c = gorm.NewDbOpener()
	cfg := &gorm.Config{
		Debug: true,
		Dialect: func(s string) g.Dialector {
			return sqlite.Open(s)
		},
		Cfg: &g.Config{},
		//https://github.com/go-gorm/gorm/issues/2875
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	}
	TestUnitOfWorkManager = uow.NewManager(&uow.Config{SupportNestedTransaction: false}, func(ctx context.Context, kind, key string) uow.TransactionalDb {
		if kind == gorm.GormDbKind {
			db, err := TestDbOpener.Open(cfg, key)
			if err != nil {
				panic(err)
			}
			return gorm2.NewTransactionDb(db)
		}
		panic(errors.New(fmt.Sprintf("can not resolve %s", key)))
	})
	TestDbProvider := GetProvider(cfg)

	TestDb = GetDb(context.Background(), TestDbProvider)
	err := AutoMigrate(nil, TestDb)
	if err != nil {
		panic(err)
	}
	TestTenantRepo = &TenantRepo{
		DbProvider: TestDbProvider,
	}
	TestGormTenantStore = NewGormTenantStore(TestTenantRepo)
	exitCode := m.Run()

	c()
	// 退出
	os.Exit(exitCode)

}

func GetProvider(cfg *gorm.Config) *gorm.DefaultDbProvider {

	ct := common.ContextCurrentTenant{}
	conn := make(data.ConnStrings, 1)
	conn.SetDefault("file::memory:?cache=shared")
	mr := common.NewMultiTenancyConnStrResolver(ct, func() common.TenantStore {
		return *TestGormTenantStore
	}, data.NewConnStrOption(conn))
	r := gorm.NewDefaultDbProvider(mr, cfg, TestUnitOfWorkManager, TestDbOpener)
	return r
}
