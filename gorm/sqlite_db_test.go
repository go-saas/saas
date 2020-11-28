package gorm

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"testing"
)

var TestDb *g.DB
var TestDbProvider *DefaultDbProvider
var close DbClean
var TenantId1 string
var TenantId2 string

func TestMain(m *testing.M) {

	TestDbProvider, close = GetProvider()

	TestDb = GetDb(context.Background(), TestDbProvider)
	err := AutoMigrate(nil, TestDb)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	close()
	// 退出
	os.Exit(exitCode)

}

func GetProvider() (*DefaultDbProvider, DbClean) {
	cfg := Config{
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
	//use memory store
	TenantId1 = uuid.New().String()
	TenantId2 = uuid.New().String()
	ts := common.NewMemoryTenantStore(
		[]common.TenantConfig{
			{ID: TenantId1, Name: "Test1"},
			{ID: TenantId2, Name: "Test2", Conn: map[string]string{
				data.Default: ":memory:?cache=shared",
			}},
		})
	conn := make(data.ConnStrings, 1)
	conn.SetDefault("file::memory:?cache=shared")
	mr := common.NewMultiTenancyConnStrResolver(ct, func() common.TenantStore {
		return ts
	}, data.ConnStrOption{
		Conn: conn,
	})
	r, close := NewDefaultDbProvider(mr, cfg)
	return r, close
}
