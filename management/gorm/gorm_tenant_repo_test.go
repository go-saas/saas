package gorm

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/management/domain"
	"github.com/goxiaoy/go-saas/management/gorm/entity"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)


func TestGormTestTenantRepo_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		t   domain.Tenant
	}
	tests := []struct {
		name   string
		args   args
	}{
		{"Test",args{
			context.Background(),domain.Tenant{
				ID:          uuid.New().String(),
				Name:        "Test",
				DisplayName: "Test",
				Region:      "Test",
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err :=TestTenantRepo.Create(tt.args.ctx,tt.args.t)
			assert.NoError(t,err)
			var tDb entity.Tenant
			err =TestDb.Model(&entity.Tenant{}).Where("id = ?",tt.args.t.ID).First(&tDb).Error
			assert.NoError(t,err)
			var dt = new(domain.Tenant)
			common.Copy(tDb,dt)
			assert.Equal(t,tt.args.t.ID,dt.ID)
			assert.Equal(t,tt.args.t.Name,dt.Name)
			assert.Equal(t,tt.args.t.DisplayName,dt.DisplayName)
			assert.Equal(t,tt.args.t.Region,dt.Region)
		})
	}
}


func TestGormTestTenantRepo_FindByIdOrName(t *testing.T) {

	//insert
	id :=  uuid.New().String()

	es := []entity.Tenant{
		{
			ID:          id,
			Name:        "Test"+id,
		},
		{
			ID:          uuid.New().String(),
			Name:        "Test2",
		},
	}
	for _, e := range es {
		err:=TestDb.Model(&entity.Tenant{}).Create(&e).Error
		assert.NoError(t,err)
	}
	e,err :=TestTenantRepo.FindByIdOrName(context.Background(),id)
	assert.NoError(t,err)
	assert.Equal(t,e.ID,es[0].ID)
	assert.Equal(t,e.Name,es[0].Name)

	e,err =TestTenantRepo.FindByIdOrName(context.Background(),"Test"+id)
	assert.NoError(t,err)
	assert.Equal(t,e.ID,es[0].ID)
	assert.Equal(t,e.Name,es[0].Name)

	e,err =TestTenantRepo.FindByIdOrName(context.Background(),uuid.New().String())
	assert.NoError(t,err)
	assert.Equal(t,true,e==nil)


}

func TestGormTestTenantRepo_GetCount(t *testing.T) {

	preCount ,err:=TestTenantRepo.GetCount(context.Background())
	assert.NoError(t,err)

	es := []entity.Tenant{
		{
			ID:           uuid.New().String(),
			Name:        "Test",
		},
		{
			ID:          uuid.New().String(),
			Name:        "Test2",
		},
	}
	for _, e := range es {
		err:=TestDb.Model(&entity.Tenant{}).Create(&e).Error
		assert.NoError(t,err)
	}
	newCount ,err :=TestTenantRepo.GetCount(context.Background())
	assert.NoError(t,err)
	assert.Equal(t,int64(2),newCount-preCount)
}

func TestGormTestTenantRepo_GetPaged(t *testing.T) {
	//get count
	var count int64
	err:=TestDb.Model(&entity.Tenant{}).Count(&count).Error
	assert.NoError(t,err)
	wg :=sync.WaitGroup{}
	wg.Add(200)
	for i:=0;i<200;i++ {
		go func(i int) {
			e:=entity.Tenant{
					ID:           uuid.New().String(),
					Name:        "Test"+strconv.Itoa(i),
				}
			err:=TestDb.Model(&entity.Tenant{}).Create(&e).Error
			assert.NoError(t,err)
			wg.Done()
		}(i)
	}
	wg.Wait()
	c,d,err:=TestTenantRepo.GetPaged(context.Background(),common.Pagination{
		Offset: 10,
		Limit:  20,
	})
	assert.NoError(t,err)
	assert.Equal(t,int64(200),c-count)
	assert.Equal(t,20,len(d))
}

func TestGormTestTenantRepo_Update(t *testing.T) {
	id :=  uuid.New().String()
	es := []entity.Tenant{
		{
			ID:           id,
			Name:        "Test",
		},
	}
	for _, e := range es {
		err:=TestDb.Model(&entity.Tenant{}).Create(&e).Error
		assert.NoError(t,err)
	}
	TestTenantRepo.Update(context.Background(),id,domain.Tenant{
		Name:        "TestNew",
		DisplayName: "",
	})
	var tenant entity.Tenant
	err:=TestDb.Model(&entity.Tenant{}).Where("id = ?",id).First(&tenant).Error
	assert.NoError(t,err)
	assert.Equal(t,"TestNew",tenant.Name)
	assert.Equal(t,id,tenant.ID)
}