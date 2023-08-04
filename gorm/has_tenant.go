//ref:https://github.com/go-gorm/gorm/blob/master/soft_delete.go

package gorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"github.com/go-saas/saas"

	"github.com/go-saas/saas/data"
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
	ct, _ := saas.FromCurrentTenant(ctx)
	at := data.FromAutoSetTenantId(ctx)
	if at {
		if ct.GetId() != t.String {
			//mismatch
			if ct.GetId() != "" {
				//only normalize in tenant side
				if !t.Valid || t.String == "" {
					//tenant want to insert self
					return clause.Expr{SQL: "?", Vars: []interface{}{ct.GetId()}}
				} else {
					//tenant want to insert others
					//force reset
					return clause.Expr{SQL: "?", Vars: []interface{}{ct.GetId()}}
				}
			}
		}
	}
	if t.Valid && t.String != "" {
		return clause.Expr{SQL: "?", Vars: []interface{}{t.String}}
	} else {
		return clause.Expr{SQL: "?", Vars: []interface{}{nil}}
	}
}

// Scan implements the Scanner interface.
func (t *HasTenant) Scan(value interface{}) error {
	return (*sql.NullString)(t).Scan(value)
}

// Value implements the driver Valuer interface.
func (t HasTenant) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.String, nil
}

func (t HasTenant) MarshalJSON() ([]byte, error) {
	if t.Valid {
		return json.Marshal(t.String)
	}
	return json.Marshal(nil)
}

func (t HasTenant) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		t.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &t.String)
	if err == nil {
		t.Valid = true
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
	t, _ := saas.FromCurrentTenant(stmt.Context)
	e := data.FromMultiTenancyDataFilter(stmt.Context)
	if !e {
		return
	}
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
		var v interface{}
		if t.GetId() == "" {
			v = nil
		} else {
			v = t.GetId()
		}
		stmt.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: sd.Field.DBName}, Value: v},
		}})
		stmt.Clauses["multi_tenancy_enabled"] = clause.Clause{}
	}
}

func (HasTenant) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{HasTenantDeleteClause{Field: f}}
}

type HasTenantDeleteClause struct {
	Field *schema.Field
}

func (sd HasTenantDeleteClause) Name() string {
	return ""
}

func (sd HasTenantDeleteClause) Build(clause.Builder) {
}

func (sd HasTenantDeleteClause) MergeClause(*clause.Clause) {
}

func (sd HasTenantDeleteClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 {
		HasTenantQueryClause(sd).ModifyStatement(stmt)
	}
}

func (HasTenant) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{HasTenantUpdateClause{Field: f}}
}

type HasTenantUpdateClause struct {
	Field *schema.Field
}

func (sd HasTenantUpdateClause) Name() string {
	return ""
}

func (sd HasTenantUpdateClause) Build(clause.Builder) {
}

func (sd HasTenantUpdateClause) MergeClause(*clause.Clause) {
}

func (sd HasTenantUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if stmt.SQL.Len() == 0 {
		HasTenantQueryClause(sd).ModifyStatement(stmt)
	}
}
