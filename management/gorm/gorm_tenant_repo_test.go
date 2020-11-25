package gorm

import (
	"context"
	"github.com/google/uuid"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"github.com/goxiaoy/go-saas/gorm"
	"github.com/goxiaoy/go-saas/management/domain"
	"gorm.io/driver/sqlite"
	g "gorm.io/gorm"
	"os"
	"reflect"
	"testing"
)

var tenantRepo GormTenantRepo

func TestMain(m *testing.M) {
	cfg := gorm.Config{
		Debug:        true,
		Dialect: func(s string) g.Dialector {
			return sqlite.Open(s)
		},
		Cfg:          nil,
	}
	ct := common.ContextCurrentTenant{}
	ts := GormTenantStore{}
	conn := make(data.ConnStrings,1)
	conn.SetDefault("file::memory:?cache=shared")
	mr := common.NewMultiTenancyConnStrResolver(ct,ts,data.ConnStrOption{
		Conn: conn,
	})
	r ,close := gorm.NewDefaultDbProvider(mr,cfg)
	tenantRepo = GormTenantRepo{
		DbProvider: r,
	}
	exitCode := m.Run()

	close()
	// 退出
	os.Exit(exitCode)

}

func TestGormTenantRepo_Create(t *testing.T) {
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
			tenantRepo.Create(tt.args.ctx,tt.args.t)
		})
	}
}

func TestGormTenantRepo_Db(t *testing.T) {
	type fields struct {
		DbProvider gorm.DbProvider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *g.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GormTenantRepo{
				DbProvider: tt.fields.DbProvider,
			}
			if got := g.Db(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Db() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormTenantRepo_FindByIdOrName(t *testing.T) {
	type fields struct {
		DbProvider gorm.DbProvider
	}
	type args struct {
		ctx      context.Context
		idOrName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *domain.Tenant
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GormTenantRepo{
				DbProvider: tt.fields.DbProvider,
			}
			if got,_ := g.FindByIdOrName(tt.args.ctx, tt.args.idOrName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByIdOrName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormTenantRepo_GetCount(t *testing.T) {
	type fields struct {
		DbProvider gorm.DbProvider
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GormTenantRepo{
				DbProvider: tt.fields.DbProvider,
			}
			if got,_ := g.GetCount(tt.args.ctx); got != tt.want {
				t.Errorf("GetCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGormTenantRepo_GetPaged(t *testing.T) {
	type fields struct {
		DbProvider gorm.DbProvider
	}
	type args struct {
		ctx context.Context
		p   common.Pagination
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantC  int64
		wantT  []*domain.Tenant
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GormTenantRepo{
				DbProvider: tt.fields.DbProvider,
			}
			gotC, gotT,_ := g.GetPaged(tt.args.ctx, tt.args.p)
			if gotC != tt.wantC {
				t.Errorf("GetPaged() gotC = %v, want %v", gotC, tt.wantC)
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("GetPaged() gotT = %v, want %v", gotT, tt.wantT)
			}
		})
	}
}

func TestGormTenantRepo_Update(t *testing.T) {
	type fields struct {
		DbProvider gorm.DbProvider
	}
	type args struct {
		ctx context.Context
		id  string
		t   domain.Tenant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}