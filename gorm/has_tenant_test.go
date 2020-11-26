package gorm

import (
	"context"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/stretchr/testify/assert"
	g "gorm.io/gorm"
	"os"
	"sync"
	"testing"
)

type TestEntity struct {
	ID string
	//custom field
	TenantId HasTenant
}

func GetDb(ctx context.Context, provider DbProvider) *g.DB {
	return provider.Get(ctx, "Test")
}
func AutoMigrate(f func(*g.DB), db *g.DB) error {
	if f != nil {
		f(db)
	}
	return db.AutoMigrate(
		new(TestEntity),
	)
}

var db *g.DB

func TestMain(m *testing.M) {

	r, close := GetProvider()

	db = GetDb(context.Background(), r)
	err := AutoMigrate(nil, db)
	if err != nil {
		panic(err)
	}

	exitCode := m.Run()
	close()
	// 退出
	os.Exit(exitCode)

}
func TestCustomField(t *testing.T) {

	//insert records
	i := []TestEntity{
		{ID: "Host1", TenantId: HasTenant{
			Valid: false,
		}},
		{ID: "Host2", TenantId: HasTenant{
			Valid: false,
		}},
		{ID: "TenantA1", TenantId: HasTenant{
			String: "A",
			Valid: true,
		}},
		{ID: "TenantA2", TenantId: HasTenant{
			String: "A",
			Valid: true,
		}},
		{ID: "TenantB1", TenantId: HasTenant{
			String: "B",
			Valid: true,
		}},
		{ID: "TenantB2", TenantId: HasTenant{
			String: "B",
			Valid: true,
		}},
	}
	wg := sync.WaitGroup{}
	wg.Add(len(i))
	for _, entity := range i {
		go func(entity TestEntity) {
			db.Create(&entity)
			wg.Done()
		}(entity)
	}
	wg.Wait()

	disableCtx := data.NewDisableMultiTenancyDataFilter(context.Background())

	var count int64
	//check count
	tx:=db.WithContext(disableCtx).Model(&TestEntity{}).Count(&count)
	assert.NoError(t,tx.Error)
	assert.Equal(t,int64(len(i)),count)


	t.Run("Host", func(t *testing.T) {
		ctx :=  common.NewCurrentTenant(context.Background(),"","")
		var count int64
		tx:= db.WithContext(ctx).Model(&TestEntity{}).Count(&count)
		assert.NoError(t,tx.Error)
		assert.Equal(t,int64(2),count)
	})

	t.Run("Tenant", func(t *testing.T) {
		{
			ctx :=  common.NewCurrentTenant(context.Background(),"A","")
			var count int64
			tx:= db.WithContext(ctx).Model(&TestEntity{}).Count(&count)
			assert.NoError(t,tx.Error)
			assert.Equal(t,int64(2),count)
		}
		{
			ctx :=  common.NewCurrentTenant(context.Background(),"B","")
			var count int64
			tx:= db.WithContext(ctx).Model(&TestEntity{}).Count(&count)
			assert.NoError(t,tx.Error)
			assert.Equal(t,int64(2),count)
		}
	})
}
