package common

import "database/sql"

//has tenant type
//TODO how to handle orm unrelated stuff
// Deprecated: can not handle orm unrelated
type HasTenant sql.NullString

