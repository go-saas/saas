package gorm

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"testing"
)

var TestDb *g.DB

var TestDbProvider DbProvider

var (
	TenantId1 = uuid.New().String()
	TenantId2 = uuid.New().String()
)

func TestMain(m *testing.M) {

	clientProvider := ClientProviderFunc(func(ctx context.Context, s string) (*g.DB, error) {
		db, err := sql.Open("sqlite3", s)
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
		return client.WithContext(ctx).Debug(), err
	})
	TestDbProvider = NewDbProvider(GetConnStrResolver(), clientProvider)

	TestDb = GetDb(context.Background(), TestDbProvider)
	err := AutoMigrate(nil, TestDb)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	NewDbWrap(TestDb).Close()
	// 退出
	os.Exit(exitCode)

}

func GetConnStrResolver() *common.MultiTenancyConnStrResolver {
	//use memory store

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
	return mr
}
