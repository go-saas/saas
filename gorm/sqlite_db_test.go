package gorm

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/uow"
	"github.com/goxiaoy/uow/gorm"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"strings"
	"testing"
)

var TestDb *g.DB
var TestDbProvider *DefaultDbProvider
var c DbClean
var TenantId1 string
var TenantId2 string
var TestDbOpener DbOpener

var TestUnitOfWorkManager uow.Manager

func TestMain(m *testing.M) {

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
	TestUnitOfWorkManager = uow.NewManager(&uow.Config{SupportNestedTransaction: false}, func(ctx context.Context, key string) uow.TransactionalDb {
		if strings.HasPrefix(key, "gorm_") {
			db, err := TestDbOpener.Open(cfg, strings.TrimLeft(key, "gorm_"))
			if err != nil {
				panic(err)
			}
			return gorm.NewTransactionDb(db)
		}
		panic(errors.New(fmt.Sprintf("can not resolve %s", key)))
	})

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
	}, data.NewConnStrOption(conn))
	r := NewDefaultDbProvider(mr, cfg, TestUnitOfWorkManager, TestDbOpener)
	return r
}
