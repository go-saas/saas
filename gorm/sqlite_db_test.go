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
var TenantId1 string
var TenantId2 string
var TestDbOpener DbOpener

func TestMain(m *testing.M) {
	var c func()
	TestDbOpener, c = NewDbOpener()
	cfg := &Config{
		Debug: true,
		Dialect: func(s string) g.Dialector {
			return sqlite.Open(s)
		},
		Cfg: &g.Config{},
		//https://github.com/go-gorm/gorm/issues/2875
		MaxOpenConn: 1,
		MaxIdleConn: 1,
	}

	TestDbProvider = GetProvider(cfg)

	TestDb = GetDb(context.Background(), TestDbProvider)
	err := AutoMigrate(nil, TestDb)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	c()
	// 退出
	os.Exit(exitCode)

}

func GetProvider(cfg *Config) *DefaultDbProvider {
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
	mr := common.NewMultiTenancyConnStrResolver(func() common.TenantStore {
		return ts
	}, data.NewConnStrOption(conn))
	p := NewDefaultDbProvider(mr, cfg, TestDbOpener)
	return p
}
