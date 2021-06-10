package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/uow"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"testing"
)

var TestDb *g.DB
var TestDbProvider *gorm.DefaultDbProvider
var TestTenantRepo *TenantRepo
var TestGormTenantStore *GormTenantStore
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
