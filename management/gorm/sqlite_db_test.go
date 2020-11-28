package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"testing"
)

var TestDb *g.DB
var TestDbProvider *gorm.DefaultDbProvider
var TestTenantRepo *GormTenantRepo
var TestGormTenantStore *GormTenantStore

func TestMain(m *testing.M) {

	TestDbProvider, close := GetProvider()

	TestDb = GetDb(context.Background(), TestDbProvider)
	err := AutoMigrate(nil, TestDb)
	if err != nil {
		panic(err)
	}
	TestTenantRepo = &GormTenantRepo{
		DbProvider: TestDbProvider,
	}
	TestGormTenantStore = NewGormTenantStore(TestTenantRepo)
	exitCode := m.Run()

	close()
	// 退出
	os.Exit(exitCode)

}

func GetProvider() (*gorm.DefaultDbProvider, gorm.DbClean) {
	cfg := gorm.Config{
		Debug: true,
		Dialect: func(s string) g.Dialector {
			return sqlite.Open(s)
		},
		Cfg: &g.Config{},
		//https://github.com/go-gorm/gorm/issues/2875
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}
	ct := common.ContextCurrentTenant{}
	conn := make(data.ConnStrings, 1)
	conn.SetDefault("file::memory:?cache=shared")
	mr := common.NewMultiTenancyConnStrResolver(ct, func() common.TenantStore {
		return *TestGormTenantStore
	}, data.ConnStrOption{
		Conn: conn,
	})
	r, close := gorm.NewDefaultDbProvider(mr, cfg)
	return r, close
}
