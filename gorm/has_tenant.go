//ref:https://github.com/go-gorm/gorm/blob/master/soft_delete.go

package gorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/goxiaoy/go-saas/common"
	"github.com/goxiaoy/go-saas/data"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type HasTenant sql.NullString

func NewTenantId(s string) HasTenant {
	if s == "" {
		return HasTenant{
			Valid: false,
		}
	} else {
		return HasTenant{
			String: s,
			Valid:  true,
		}
	}
}

func (t HasTenant) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
	ct := common.FromCurrentTenant(ctx)
	if ct.Id != t.String {
		//mismatch
		at := data.FromAutoSetTenantId(ctx)
		if at && ct.Id != "" {
			if !t.Valid || t.String == "" {
				//tenant want to insert self
				return clause.Expr{SQL: "?", Vars: []interface{}{ct.Id}}
			} else {
				//tenant wnt to insert others
				//force reset
				return clause.Expr{SQL: "?", Vars: []interface{}{ct.Id}}
			}
		}
	}
	return clause.Expr{SQL: "?", Vars: []interface{}{t}}
}

// Scan implements the Scanner interface.
func (n *HasTenant) Scan(value interface{}) error {
	return (*sql.NullString)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n HasTenant) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

func (n HasTenant) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n *HasTenant) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.String)
	if err == nil {
		n.Valid = true
	}
	return err
}

func (HasTenant) QueryClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{HasTenantQueryClause{Field: f}}
}

type HasTenantQueryClause struct {
	Field *schema.Field
}

func (sd HasTenantQueryClause) Name() string {
	return ""
}

func (sd HasTenantQueryClause) Build(clause.Builder) {
}

func (sd HasTenantQueryClause) MergeClause(*clause.Clause) {
}

func (sd HasTenantQueryClause) ModifyStatement(stmt *gorm.Statement) {
	t := common.FromCurrentTenant(stmt.Context)
	e := data.FromMultiTenancyDataFilter(stmt.Context)
	if _, ok := stmt.Clauses["multi_tenancy_enabled"]; !ok {
		if c, ok := stmt.Clauses["WHERE"]; ok {
			if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
				for _, expr := range where.Exprs {
					if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
						where.Exprs = []clause.Expression{clause.And(where.Exprs...)}
						c.Expression = where
						stmt.Clauses["WHERE"] = c
						break
					}
				}
			}
		}
		if e {
			var v interface{}
			if t.Id == "" {
				v = nil
			} else {
				v = t.Id
			}
			stmt.AddClause(clause.Where{Exprs: []clause.Expression{
				clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: v},
			}})
		}
		stmt.Clauses["multi_tenancy_enabled"] = clause.Clause{}
	}
}
