package saas

type MultiTenancySide int32

const (
	Tenant MultiTenancySide = 1 << 0
	Host   MultiTenancySide = 1 << 1
	Both                    = Tenant | Host
)
