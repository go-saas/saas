package gorm

import (
	"context"
	"github.com/google/uuid"
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
		{ID: "Host1", MultiTenancy: MultiTenancy{NewTenantId("")} },
		{ID: "Host2", MultiTenancy: MultiTenancy{ NewTenantId("")}},
		{ID: "TenantA1", MultiTenancy: MultiTenancy{ NewTenantId("A")}},
		{ID: "TenantA2", MultiTenancy: MultiTenancy{ NewTenantId("A")}},
		{ID: "TenantB1", MultiTenancy: MultiTenancy{NewTenantId("B")}},
		{ID: "TenantB2", MultiTenancy: MultiTenancy{ NewTenantId("B")}},
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

func TestAutoSetTenant(t *testing.T) {

	t.Run("HostAutoSet", func(t *testing.T) {
		ctx:=context.Background()
		i:= TestEntity{ID: "HostAutoSetTenant1", MultiTenancy: MultiTenancy{NewTenantId("")}}
		err :=db.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t,err)
		//find
		var hostAutoSet TestEntity
		err =db.WithContext(ctx).Model(&TestEntity{}).Where("id = ?","HostAutoSetTenant1").First(&hostAutoSet).Error
		assert.NoError(t,err)
		assert.Equal(t,hostAutoSet.TenantId,NewTenantId(""))

		i= TestEntity{ID: "HostAutoSetTenant2", MultiTenancy: MultiTenancy{NewTenantId(uuid.New().String())}}

		err =db.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t,err)

		//can not find
		hostAutoSet = TestEntity{}
		err =db.WithContext(ctx).Model(&TestEntity{}).Where("id = ?","HostAutoSetTenant2").First(&hostAutoSet).Error
		assert.Error(t,err,g.ErrRecordNotFound)

		//can find now
		hostAutoSet = TestEntity{}
		disableCtx:= data.NewDisableMultiTenancyDataFilter(ctx)
		err =db.WithContext(disableCtx).Model(&TestEntity{}).Where("id = ?","HostAutoSetTenant2").First(&hostAutoSet).Error
		assert.NoError(t,err)

	})


	t.Run("TenantAutoSet", func(t *testing.T) {
		tenantId := uuid.New().String()
		tenantId2 := uuid.New().String()
		ctx:=common.NewCurrentTenant(context.Background(),tenantId,"")
		i:= TestEntity{ID: "TenantAutoSetTenant1", MultiTenancy: MultiTenancy{NewTenantId("")}}
		err :=db.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t,err)
		//find
		var TenantAutoSet TestEntity
		err =db.WithContext(ctx).Model(&TestEntity{}).Where("id = ?","TenantAutoSetTenant1").First(&TenantAutoSet).Error
		assert.NoError(t,err)
		assert.Equal(t,TenantAutoSet.TenantId.String,tenantId)

		i= TestEntity{ID: "TenantAutoSetTenant2", MultiTenancy: MultiTenancy{NewTenantId(tenantId2)}}

		err =db.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t,err)

		//can  find
		TenantAutoSet = TestEntity{}
		err =db.WithContext(ctx).Model(&TestEntity{}).Where("id = ?","TenantAutoSetTenant2").First(&TenantAutoSet).Error
		assert.NoError(t,err)
		assert.Equal(t,TenantAutoSet.TenantId.String,tenantId)

		//disable auto set
		ctx = data.NewDisableAutoSetTenantId(ctx)

		i= TestEntity{ID: "TenantAutoSetTenant3", MultiTenancy: MultiTenancy{NewTenantId(tenantId2)}}
		err =db.WithContext(ctx).Model(&TestEntity{}).Create(&i).Error
		assert.NoError(t,err)
		
		//can not find
		TenantAutoSet = TestEntity{}
		err =db.WithContext(ctx).Model(&TestEntity{}).Where("id = ?","TenantAutoSetTenant3").First(&TenantAutoSet).Error
		assert.Error(t,err)

		//can find now
		TenantAutoSet = TestEntity{}
		disableCtx:= data.NewDisableMultiTenancyDataFilter(ctx)
		err =db.WithContext(disableCtx).Model(&TestEntity{}).Where("id = ?","TenantAutoSetTenant3").First(&TenantAutoSet).Error
		assert.NoError(t,err)
		assert.Equal(t,TenantAutoSet.TenantId.String,tenantId2)

	})




}
