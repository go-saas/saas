package main

import (
	"context"
	"fmt"
	"github.com/goxiaoy/go-saas"

	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/seed"
	gorm2 "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Seed struct {
	dbProvider gorm.DbProvider
	connStrGen saas.ConnStrGenerator
}

func NewSeed(dbProvider gorm.DbProvider, connStrGen saas.ConnStrGenerator) *Seed {
	return &Seed{dbProvider: dbProvider, connStrGen: connStrGen}
}

func (s *Seed) Seed(ctx context.Context, sCtx *seed.Context) error {
	db := s.dbProvider.Get(ctx, "")

	if sCtx.TenantId == "" {
		//seed host
		t3 := Tenant{ID: "3", Name: "Test3"}
		t3Conn, _ := s.connStrGen.Gen(ctx, saas.NewBasicTenantInfo(t3.ID, t3.Name))
		t3.Conn = []TenantConn{
			{Key: data.Default, Value: t3Conn}, // use tenant3.db
		}
		err := db.Model(&Tenant{}).Clauses(clause.OnConflict{UpdateAll: true}).CreateInBatches([]Tenant{
			{ID: "1", Name: "Test1"}, // use default shared.db
			{ID: "2", Name: "Test2"},
			t3}, 10).Error
		if err != nil {
			return err
		}
		entities := []Post{
			{
				Model:       gorm2.Model{ID: 1},
				Title:       fmt.Sprintf("Host Side"),
				Description: fmt.Sprintf("Hello Host"),
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
			},
			{
				Model:       gorm2.Model{ID: 4},
				Title:       fmt.Sprintf("Tenant %s Post 2", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are two posts in this tenant. This is post 2", sCtx.TenantId),
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
			},
			{
				Model:       gorm2.Model{ID: 6},
				Title:       fmt.Sprintf("Tenant %s Post 2", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are there posts in this tenant. This is post 2", sCtx.TenantId),
			},
			{
				Model:       gorm2.Model{ID: 7},
				Title:       fmt.Sprintf("Tenant %s Post 2", sCtx.TenantId),
				Description: fmt.Sprintf("Hello from tenant %s. There are there posts in this tenant. This is post 3", sCtx.TenantId),
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
