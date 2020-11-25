package data

import (
	"context"
)

type DefaultConnStrResolver struct {
	Opt ConnStrOption
}

//direct return value from option value
func (d DefaultConnStrResolver) Resolve(_ context.Context, key string) string {
	if key!=""{
		v:=d.Opt.Conn[key]
		if v!=""{
			return v
		}
	}
	return d.Opt.Conn.Default()
}



