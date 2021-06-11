package gorm

import (
	"gorm.io/gorm"
	"sync"
	"time"
)

type DbClean func()

type DbOpener interface {
	Open(c *Config, s string) (*gorm.DB, error)
}

type dbOpener struct {
	mtx sync.Mutex
	db  map[string]*gorm.DB
}

func NewDbOpener() (DbOpener, DbClean) {
	m := make(map[string]*gorm.DB)
	close := func() {
		for _, d := range m {
			closeDb(d)
		}
	}
	return &dbOpener{
		db: m,
	}, close
}

func (d *dbOpener) Open(c *Config, s string) (*gorm.DB, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	db, ok := d.db[s]
	if ok {
		return db, nil
	}

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

	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)

	d.db[s] = db
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