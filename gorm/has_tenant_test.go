package gorm

import (
	"context"
	"github.com/go-saas/saas"
	"github.com/google/uuid"

	"github.com/go-saas/saas/data"
	"github.com/stretchr/testify/assert"
	g "gorm.io/gorm"
	"sync"
	"testing"
)

type TestEntity struct {
	ID string
	MultiTenancy
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

func TestCustomField(t *testing.T) {

	//insert records
	i := []TestEntity{
		{ID: "Host1", MultiTenancy: MultiTenancy{NewTenantId("")}},
		{ID: "Host2", MultiTenancy: MultiTenancy{NewTenantId("")}},
		{ID: "TenantA1", MultiTenancy: MultiTenancy{NewTenantId("A")}},
		{ID: "TenantA2", MultiTenancy: MultiTenancy{NewTenantId("A")}},
		{ID: "TenantB1", MultiTenancy: MultiTenancy{NewTenantId("B")}},
		{ID: "TenantB2", MultiTenancy: MultiTenancy{NewTenantId("B")}},
	}
	wg := sync.WaitGroup{}
	wg.Add(len(i))
	for _, entity := range i {
		go func(entity TestEntity) {
			TestDb.Create(&entity)
			wg.Done()
		}(entity)
	}
	wg.Wait()

	disableCtx := data.NewMultiTenancyDataFilter(context.Background(), false)

	var count int64
	//check count
	tx := TestDb.WithContext(disableCtx).Model(&TestEntity{}).Count(&count)
	assert.NoError(t, tx.Error)
	assert.Equal(t, int64(len(i)), count)

	t.Run("Host", func(t *testing.T) {
		ctx := saas.NewCurrentTenant(context.Background(), "", "")
		var count int64
		tx := TestDb.WithContext(ctx).Model(&TestEntity{}).Count(&count)
		assert.NoError(t, tx.Error)
		assert.Equal(t, int64(2), count)
	})

	t.Run("Tenant", func(t *testing.T) {
		{
			ctx := saas.NewCurrentTenant(context.Background(), "A", "")
			var count int64
			tx := TestDb.WithContext(ctx).Model(&TestEntity{}).Count(&count)
			assert.NoError(t, tx.Error)
			assert.Equal(t, int64(2), count)
		}
		{
			ctx := saas.NewCurrentTenant(context.Background(), "B", "")
			var count int64
			tx := TestDb.WithContext(ctx).Model(&TestEntity{}).Count(&count)
			assert.NoError(t, tx.Error)
			assert.Equal(t, int64(2), count)
		}
	})
}

func TestAutoSetTenant(t *testing.T) {

	t.Run("HostAutoSet", func(t *testing.T) {
		ctx := context.Background()
		i := TestEntity{ID: "HostAutoSetTenant1", MultiTenancy: MultiTenancy{NewTenantId("")}}
		err := TestDb.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t, err)
		//find
		var hostAutoSet TestEntity
		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Where("id = ?", "HostAutoSetTenant1").First(&hostAutoSet).Error
		assert.NoError(t, err)
		assert.Equal(t, hostAutoSet.TenantId, NewTenantId(""))

		i = TestEntity{ID: "HostAutoSetTenant2", MultiTenancy: MultiTenancy{NewTenantId(uuid.New().String())}}

		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t, err)

		//can not find
		hostAutoSet = TestEntity{}
		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Where("id = ?", "HostAutoSetTenant2").First(&hostAutoSet).Error
		assert.Error(t, err, g.ErrRecordNotFound)

		//can find now
		hostAutoSet = TestEntity{}
		disableCtx := data.NewMultiTenancyDataFilter(ctx, false)
		err = TestDb.WithContext(disableCtx).Model(&TestEntity{}).Where("id = ?", "HostAutoSetTenant2").First(&hostAutoSet).Error
		assert.NoError(t, err)

	})

	t.Run("TenantAutoSet", func(t *testing.T) {
		tenantId := uuid.New().String()
		tenantId2 := uuid.New().String()
		ctx := saas.NewCurrentTenant(context.Background(), tenantId, "")
		i := TestEntity{ID: "TenantAutoSetTenant1", MultiTenancy: MultiTenancy{NewTenantId("")}}
		err := TestDb.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t, err)
		//find
		var TenantAutoSet TestEntity
		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Where("id = ?", "TenantAutoSetTenant1").First(&TenantAutoSet).Error
		assert.NoError(t, err)
		assert.Equal(t, TenantAutoSet.TenantId.String, tenantId)

		i = TestEntity{ID: "TenantAutoSetTenant2", MultiTenancy: MultiTenancy{NewTenantId(tenantId2)}}

		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t, err)

		//can  find
		TenantAutoSet = TestEntity{}
		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Where("id = ?", "TenantAutoSetTenant2").First(&TenantAutoSet).Error
		assert.NoError(t, err)
		assert.Equal(t, TenantAutoSet.TenantId.String, tenantId)

		//disable auto set
		ctx = data.NewAutoSetTenantId(ctx, false)

		i = TestEntity{ID: "TenantAutoSetTenant3", MultiTenancy: MultiTenancy{NewTenantId(tenantId2)}}
		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t, err)

		//can not find
		TenantAutoSet = TestEntity{}
		err = TestDb.WithContext(ctx).Model(&TestEntity{}).Where("id = ?", "TenantAutoSetTenant3").First(&TenantAutoSet).Error
		assert.Error(t, err)

		//can find now
		TenantAutoSet = TestEntity{}
		disableCtx := data.NewMultiTenancyDataFilter(ctx, false)
		err = TestDb.WithContext(disableCtx).Model(&TestEntity{}).Where("id = ?", "TenantAutoSetTenant3").First(&TenantAutoSet).Error
		assert.NoError(t, err)
		assert.Equal(t, TenantAutoSet.TenantId.String, tenantId2)

	})

}
