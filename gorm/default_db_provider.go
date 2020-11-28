package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/data"
	"gorm.io/gorm"
	"sync"
	"sync/atomic"
	"time"
)

type DefaultDbProvider struct {
	//TODO performance issue
	m             sync.Map
	cs            data.ConnStrResolver
	c             Config
	createCounter uint32
}

type DbClean func()

//create db provider and close function
func NewDefaultDbProvider(cs data.ConnStrResolver, c Config) (d *DefaultDbProvider, close DbClean) {
	var m sync.Map
	close = func() {
		m.Range(func(key, value interface{}) bool {
			d, ok := value.(*gorm.DB)
			if ok {
				closeDb(d)
			}
			return true
		})
	}
	d = &DefaultDbProvider{
		m:  m,
		cs: cs,
		c:  c,
	}
	return
}

func NewDB(c *Config, s string) (*gorm.DB, error) {
	db, err := gorm.Open(c.Dialect(s), c.Cfg)
	if err != nil {
		//error
		return nil, err
	}
	if c.Debug {
		db = db.Debug()
	}
	sqlDB, err := db.DB()
	if err != nil {
		//error
		return db, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	//register callback
	db.Callback().Create().Before("gorm:create").Register(MultiTenantBeforeCreateName, AutoSetTenant)
	db.Callback().Query().Register(MultiTenantQueryName, AutoFilterCurrentTenant)

	return db, nil
}

func closeDb(d *gorm.DB) error {
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

func (d *DefaultDbProvider) Get(ctx context.Context, key string) *gorm.DB {
	//resolve connection string
	s := d.cs.Resolve(ctx, key)
	//try get db object from map
	var g *gorm.DB
	gv, ok := d.m.Load(s)
	if !ok {
		//not found
		newDb, err := NewDB(&d.c, s)
		atomic.AddUint32(&d.createCounter, 1)
		if err != nil {
			//
			panic(err)
		}
		gv, ok = d.m.LoadOrStore(s, newDb)
		if ok {
			//indicate loaded, should close newDb
			//TODO performance issue
			closeDb(newDb)
		}
		g, _ = gv.(*gorm.DB)
	} else {
		//in this map
		g, _ = gv.(*gorm.DB)

	}
	return g.WithContext(ctx)
}
