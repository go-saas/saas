package domain

import "time"

//TenantConn connection string info
type TenantConn struct {
	TenantId string
	//key of connection string
	Key string
	//connection string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
