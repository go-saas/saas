package common

import (
	"context"
	"reflect"
	"testing"
)

func TestContextCurrentTenant_Change(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   string
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Change", args{ctx: context.Background(), id: "1", name: "Test1"}},
	}
	for _, tt := range tests {
		changeSlice := []BasicTenantInfo{{}, *NewBasicTenantInfo("1", "Test1"), *NewBasicTenantInfo("2", "Test2"), {}}
		cancelSlice := make([]CancelFunc, 3)
		t.Run(tt.name, func(t *testing.T) {
			c := ContextCurrentTenant{}
			//empty head
			currentContext := tt.args.ctx
			for i := 0; i < len(changeSlice)-1; i++ {
				//change
				newCtx, cancel := c.Change(currentContext, changeSlice[i].Id, changeSlice[i].Name)
				currentContext = newCtx
				cancelSlice[i] = cancel
				if !reflect.DeepEqual(getCurrent(newCtx), changeSlice[i]) {
					t.Errorf("Change() got = %v, want %v", getCurrent(newCtx), changeSlice[i])
				}
			}
			//cancel
			for i := len(cancelSlice) - 1; i > 0; i-- {
				//cancel
				newCtx := cancelSlice[i](currentContext)
				if !reflect.DeepEqual(getCurrent(newCtx), changeSlice[i-1]) {
					t.Errorf("Change() cancel got = %v, want %v", getCurrent(newCtx), changeSlice[i])
				}
			}
			newCtx := cancelSlice[0](currentContext)
			if !reflect.DeepEqual(getCurrent(newCtx), getCurrent(tt.args.ctx)) {
				t.Errorf("Change() cancel got = %v, want %v", getCurrent(newCtx), getCurrent(tt.args.ctx))
			}
		})
	}
}

func TestContextCurrentTenant_Id(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"HostGetCurrent", args{ctx: context.Background()}, ""},
		{"TenantGetCurrent", args{ctx: NewCurrentTenant(context.Background(), "1", "Test")}, "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ContextCurrentTenant{}
			if got := c.Id(tt.args.ctx); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContextCurrentTenant_IsAvailable(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"HostGetCurrent", args{ctx: context.Background()}, false},
		{"TenantGetCurrent", args{ctx: NewCurrentTenant(context.Background(), "1", "Test")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ContextCurrentTenant{}
			if got := c.IsAvailable(tt.args.ctx); got != tt.want {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCurrent(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want BasicTenantInfo
	}{
		{"HostGetCurrent", args{ctx: context.Background()}, BasicTenantInfo{}},
		{"TenantGetCurrent", args{ctx: NewCurrentTenant(context.Background(), "1", "Test")}, *NewBasicTenantInfo("1", "Test")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCurrent(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCurrent() = %v, want %v", got, tt.want)
			}
		})
	}
}
