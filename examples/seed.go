package main

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/seed"
	"gorm.io/driver/sqlite"
	gorm2 "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Seed struct {
	dbProvider gorm.DbProvider
}

func NewSeed(dbProvider gorm.DbProvider) *Seed {
	return &Seed{dbProvider: dbProvider}
}

func (s *Seed) Seed(ctx context.Context, sCtx *seed.Context) error {
	db := s.dbProvider.Get(ctx, "")
	liteDial := db.Dialector.(*sqlite.Dialector)
	dsn := liteDial.DSN
	if sCtx.TenantId == "" {
		entities := []Post{
			{
				Model:       gorm2.Model{ID: 1},
				Title:       fmt.Sprintf("Host Side"),
				Description: fmt.Sprintf("Hello Host"),
				DSN:         dsn,
			},
		}
		if err := createPosts(db, entities); err != nil {
			return err
		}
	}
	if sCtx.TenantId == "1" {
		entities := []Post{
			{
				Model:       gorm2.Model{ID: 2},
				Title:       fmt.Sprintf("Tenant %s Post 1", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are one post in this tenant. This is post 1", sCtx.TenantId),
				DSN:         dsn,
			},
		}
		if err := createPosts(db, entities); err != nil {
			return err
		}
	}
	if sCtx.TenantId == "2" {
		entities := []Post{
			{
				Model:       gorm2.Model{ID: 3},
				Title:       fmt.Sprintf("Tenant %s Post 1", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are two posts in this tenant. This is post 1", sCtx.TenantId),
				DSN:         dsn,
			},
			{
				Model:       gorm2.Model{ID: 4},
				Title:       fmt.Sprintf("Tenant %s Post 2", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are two posts in this tenant. This is post 2", sCtx.TenantId),
				DSN:         dsn,
			},
		}
		if err := createPosts(db, entities); err != nil {
			return err
		}
	}

	if sCtx.TenantId == "3" {
		entities := []Post{
			{
				Model:       gorm2.Model{ID: 5},
				Title:       fmt.Sprintf("Tenant %s Post 1", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are there posts in this tenant. This is post 1", sCtx.TenantId),
				DSN:         dsn,
			},
			{
				Model:       gorm2.Model{ID: 6},
				Title:       fmt.Sprintf("Tenant %s Post 2", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are there posts in this tenant. This is post 2", sCtx.TenantId),
				DSN:         dsn,
			},
			{
				Model:       gorm2.Model{ID: 7},
				Title:       fmt.Sprintf("Tenant %s Post 2", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are there posts in this tenant. This is post 3", sCtx.TenantId),
				DSN:         dsn,
			},
		}
		if err := createPosts(db, entities); err != nil {
			return err
		}
	}
	return nil
}

func createPosts(db *gorm2.DB, entities []Post) error {
	for _, entity := range entities {
		err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Model(&Post{}).Create(&entity).Error
		if err != nil {
			return err
		}
	}
	return nil
}
