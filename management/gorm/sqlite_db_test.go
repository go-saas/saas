package gorm

import (
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
)

func GetProvider() (*gorm.DefaultDbProvider, gorm.DbClean) {
	cfg := gorm.Config{
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
	ts := GormTenantStore{}
	conn := make(data.ConnStrings,1)
	conn.SetDefault("file::memory:?cache=shared")
	mr := common.NewMultiTenancyConnStrResolver(ct,ts,data.ConnStrOption{
		Conn: conn,
	})
	r ,close := gorm.NewDefaultDbProvider(mr,cfg)
	return r,close
}
