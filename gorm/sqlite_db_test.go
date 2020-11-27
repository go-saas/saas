package gorm

import (
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
)

func GetProvider() (*DefaultDbProvider, DbClean) {
	cfg := Config{
		Debug:        true,
		Dialect: func(s string) g.Dialector {
			return sqlite.Open(s)
		},
		Cfg:&g.Config{
		},
		//https://github.com/go-gorm/gorm/issues/2875
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	}
	ct := common.ContextCurrentTenant{}
	//use memory store
	ts := common.NewMemoryTenantStore(
		[]common.TenantConfig{
			{ID: "1",Name: "Test1"},
			{ID: "2",Name: "Test3"},
		})
	conn := make(data.ConnStrings,1)
	conn.SetDefault("file::memory:?cache=shared")
	mr := common.NewMultiTenancyConnStrResolver(ct, func() common.TenantStore {
		return ts
	},data.ConnStrOption{
		Conn: conn,
	})
	r ,close := NewDefaultDbProvider(mr,cfg)
	return r,close
}