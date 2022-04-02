package gorm

import (
	"gorm.io/gorm"
	"sync"
	"time"
)

// DbOpener open a real *gorm.DB instance by using connection string
type DbOpener interface {
	Open(c *Config, s string) (*gorm.DB, error)
	Close()
}

type dbOpener struct {
	mtx           sync.Mutex
	db            map[string]*gorm.DB
	creationHooks []DbCreationHook
}

type DbCreationHook func(db *gorm.DB) *gorm.DB

func NewDbOpener(creationHooks ...DbCreationHook) (DbOpener, func()) {
	ret := &dbOpener{
		db:            make(map[string]*gorm.DB),
		creationHooks: creationHooks,
	}
	return ret, ret.Close
}

func (d *dbOpener) Open(c *Config, s string) (*gorm.DB, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	db, ok := d.db[s]
	if ok {
		return db, nil
	}
	if c.EnsureDbExist != nil {
		if err := c.EnsureDbExist(c, s); err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(c.Dialect(s), c.Cfg)
	if err != nil {
		//error
		return nil, err
	}
	for _, ch := range d.creationHooks {
		db = ch(db)
	}

	if c.Debug {
		db = db.Debug()
	}
	sqlDB, err := db.DB()
	if err != nil {
		//error
		return db, err
	}
	if c.MaxIdleConn != nil {
		sqlDB.SetMaxIdleConns(*c.MaxIdleConn)
	}
	if c.MaxOpenConn != nil {
		sqlDB.SetMaxOpenConns(*c.MaxOpenConn)
	}
	if c.MaxLifetime != nil {
		sqlDB.SetConnMaxLifetime(time.Duration(*c.MaxLifetime) * time.Second)
	}
	d.db[s] = db
	return db, nil
}

func (d *dbOpener) Close() {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	for _, d := range d.db {
		closeDb(d)
	}

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
