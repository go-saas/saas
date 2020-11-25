package domain

import "time"

type Tenant struct {
	// unique id
	ID string
	//unique name. usually for domain name
	Name string
	//localed display name
	DisplayName string
	//region of this tenant. Useful for data storage location
	Region string
	CreatedAt time.Time
	UpdatedAt time.Time
	//should apply soft delete
	DeletedAt *time.Time
	//connection
	Conn []TenantConn
	//edition
	Features []TenantFeature
}