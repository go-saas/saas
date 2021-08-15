package gorm

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/management/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGormTenantStore_GetByNameOrId(t *testing.T) {
	//insert
	id := uuid.New().String()
	dt := domain.Tenant{
		ID:   id,
		Name: "Test",
		Conn: []domain.TenantConn{
			{TenantId: id, Key: data.Default, Value: "A"},
			{TenantId: id, Key: "B", Value: "B"},
		},
		Features: nil,
	}
	ctx := context.Background()
	err := TestTenantRepo.Create(ctx, dt)
	assert.NoError(t, err)
	tc, err := TestGormTenantStore.GetByNameOrId(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, tc.ID)
	assert.Equal(t, "Test", tc.Name)
	assert.Equal(t, "A", tc.Conn.Default())
	assert.Equal(t, "A", tc.Conn.GetOrDefault("Nil"))
	assert.Equal(t, "B", tc.Conn.GetOrDefault("B"))
	// change to tenant A
	ctx = common.NewCurrentTenant(ctx, id, "Test")

	tc, err = TestGormTenantStore.GetByNameOrId(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, tc.ID)
	assert.Equal(t, "Test", tc.Name)
	assert.Equal(t, "A", tc.Conn.Default())
	assert.Equal(t, "A", tc.Conn.GetOrDefault("Nil"))
	assert.Equal(t, "B", tc.Conn.GetOrDefault("B"))

}
