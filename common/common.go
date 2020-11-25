package common

import (
	"github.com/google/wire"
)

var CommonSet=wire.NewSet(
	DefaultMultiTenancyOption,
	wire.Struct(new(TenantResolveOption),"*"),
	wire.Struct(new(DefaultTenantResolver),"*"),
	wire.Bind(new(TenantResolver),new(*DefaultTenantResolver)),
	wire.Struct(new(MemoryTenantStore),"*"),
	wire.Bind(new(TenantStore),new(*MemoryTenantStore)),
	wire.Struct(new(DefaultTenantConfigProvider),"*"),
	wire.Bind(new(TenantConfigProvider),new(*DefaultTenantConfigProvider)),
	wire.Struct(new(ContextCurrentTenant)),
	wire.Bind(new(CurrentTenant),new(*ContextCurrentTenant)),
)
