// Code generated by entc, DO NOT EDIT.

package ent

import (
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent/post"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent/predicate"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent/tenant"
	"github.com/goxiaoy/go-saas/examples/ent/shared/ent/tenantconn"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/entql"
	"entgo.io/ent/schema/field"
)

// schemaGraph holds a representation of ent/schema at runtime.
var schemaGraph = func() *sqlgraph.Schema {
	graph := &sqlgraph.Schema{Nodes: make([]*sqlgraph.Node, 3)}
	graph.Nodes[0] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   post.Table,
			Columns: post.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: post.FieldID,
			},
		},
		Type: "Post",
		Fields: map[string]*sqlgraph.FieldSpec{
			post.FieldTenantID:    {Type: field.TypeString, Column: post.FieldTenantID},
			post.FieldTitle:       {Type: field.TypeString, Column: post.FieldTitle},
			post.FieldDescription: {Type: field.TypeString, Column: post.FieldDescription},
			post.FieldDsn:         {Type: field.TypeString, Column: post.FieldDsn},
		},
	}
	graph.Nodes[1] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   tenant.Table,
			Columns: tenant.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tenant.FieldID,
			},
		},
		Type: "Tenant",
		Fields: map[string]*sqlgraph.FieldSpec{
			tenant.FieldCreateTime:  {Type: field.TypeTime, Column: tenant.FieldCreateTime},
			tenant.FieldUpdateTime:  {Type: field.TypeTime, Column: tenant.FieldUpdateTime},
			tenant.FieldName:        {Type: field.TypeString, Column: tenant.FieldName},
			tenant.FieldDisplayName: {Type: field.TypeString, Column: tenant.FieldDisplayName},
			tenant.FieldRegion:      {Type: field.TypeString, Column: tenant.FieldRegion},
		},
	}
	graph.Nodes[2] = &sqlgraph.Node{
		NodeSpec: sqlgraph.NodeSpec{
			Table:   tenantconn.Table,
			Columns: tenantconn.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tenantconn.FieldID,
			},
		},
		Type: "TenantConn",
		Fields: map[string]*sqlgraph.FieldSpec{
			tenantconn.FieldCreateTime: {Type: field.TypeTime, Column: tenantconn.FieldCreateTime},
			tenantconn.FieldUpdateTime: {Type: field.TypeTime, Column: tenantconn.FieldUpdateTime},
			tenantconn.FieldKey:        {Type: field.TypeString, Column: tenantconn.FieldKey},
			tenantconn.FieldValue:      {Type: field.TypeString, Column: tenantconn.FieldValue},
		},
	}
	graph.MustAddE(
		"conn",
		&sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tenant.ConnTable,
			Columns: []string{tenant.ConnColumn},
			Bidi:    false,
		},
		"Tenant",
		"TenantConn",
	)
	return graph
}()

// predicateAdder wraps the addPredicate method.
// All update, update-one and query builders implement this interface.
type predicateAdder interface {
	addPredicate(func(s *sql.Selector))
}

// addPredicate implements the predicateAdder interface.
func (pq *PostQuery) addPredicate(pred func(s *sql.Selector)) {
	pq.predicates = append(pq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the PostQuery builder.
func (pq *PostQuery) Filter() *PostFilter {
	return &PostFilter{pq}
}

// addPredicate implements the predicateAdder interface.
func (m *PostMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the PostMutation builder.
func (m *PostMutation) Filter() *PostFilter {
	return &PostFilter{m}
}

// PostFilter provides a generic filtering capability at runtime for PostQuery.
type PostFilter struct {
	predicateAdder
}

// Where applies the entql predicate on the query filter.
func (f *PostFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[0].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *PostFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(post.FieldID))
}

// WhereTenantID applies the entql string predicate on the tenant_id field.
func (f *PostFilter) WhereTenantID(p entql.StringP) {
	f.Where(p.Field(post.FieldTenantID))
}

// WhereTitle applies the entql string predicate on the title field.
func (f *PostFilter) WhereTitle(p entql.StringP) {
	f.Where(p.Field(post.FieldTitle))
}

// WhereDescription applies the entql string predicate on the description field.
func (f *PostFilter) WhereDescription(p entql.StringP) {
	f.Where(p.Field(post.FieldDescription))
}

// WhereDsn applies the entql string predicate on the dsn field.
func (f *PostFilter) WhereDsn(p entql.StringP) {
	f.Where(p.Field(post.FieldDsn))
}

// addPredicate implements the predicateAdder interface.
func (tq *TenantQuery) addPredicate(pred func(s *sql.Selector)) {
	tq.predicates = append(tq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the TenantQuery builder.
func (tq *TenantQuery) Filter() *TenantFilter {
	return &TenantFilter{tq}
}

// addPredicate implements the predicateAdder interface.
func (m *TenantMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the TenantMutation builder.
func (m *TenantMutation) Filter() *TenantFilter {
	return &TenantFilter{m}
}

// TenantFilter provides a generic filtering capability at runtime for TenantQuery.
type TenantFilter struct {
	predicateAdder
}

// Where applies the entql predicate on the query filter.
func (f *TenantFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[1].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *TenantFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(tenant.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *TenantFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(tenant.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *TenantFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(tenant.FieldUpdateTime))
}

// WhereName applies the entql string predicate on the name field.
func (f *TenantFilter) WhereName(p entql.StringP) {
	f.Where(p.Field(tenant.FieldName))
}

// WhereDisplayName applies the entql string predicate on the display_name field.
func (f *TenantFilter) WhereDisplayName(p entql.StringP) {
	f.Where(p.Field(tenant.FieldDisplayName))
}

// WhereRegion applies the entql string predicate on the region field.
func (f *TenantFilter) WhereRegion(p entql.StringP) {
	f.Where(p.Field(tenant.FieldRegion))
}

// WhereHasConn applies a predicate to check if query has an edge conn.
func (f *TenantFilter) WhereHasConn() {
	f.Where(entql.HasEdge("conn"))
}

// WhereHasConnWith applies a predicate to check if query has an edge conn with a given conditions (other predicates).
func (f *TenantFilter) WhereHasConnWith(preds ...predicate.TenantConn) {
	f.Where(entql.HasEdgeWith("conn", sqlgraph.WrapFunc(func(s *sql.Selector) {
		for _, p := range preds {
			p(s)
		}
	})))
}

// addPredicate implements the predicateAdder interface.
func (tcq *TenantConnQuery) addPredicate(pred func(s *sql.Selector)) {
	tcq.predicates = append(tcq.predicates, pred)
}

// Filter returns a Filter implementation to apply filters on the TenantConnQuery builder.
func (tcq *TenantConnQuery) Filter() *TenantConnFilter {
	return &TenantConnFilter{tcq}
}

// addPredicate implements the predicateAdder interface.
func (m *TenantConnMutation) addPredicate(pred func(s *sql.Selector)) {
	m.predicates = append(m.predicates, pred)
}

// Filter returns an entql.Where implementation to apply filters on the TenantConnMutation builder.
func (m *TenantConnMutation) Filter() *TenantConnFilter {
	return &TenantConnFilter{m}
}

// TenantConnFilter provides a generic filtering capability at runtime for TenantConnQuery.
type TenantConnFilter struct {
	predicateAdder
}

// Where applies the entql predicate on the query filter.
func (f *TenantConnFilter) Where(p entql.P) {
	f.addPredicate(func(s *sql.Selector) {
		if err := schemaGraph.EvalP(schemaGraph.Nodes[2].Type, p, s); err != nil {
			s.AddError(err)
		}
	})
}

// WhereID applies the entql int predicate on the id field.
func (f *TenantConnFilter) WhereID(p entql.IntP) {
	f.Where(p.Field(tenantconn.FieldID))
}

// WhereCreateTime applies the entql time.Time predicate on the create_time field.
func (f *TenantConnFilter) WhereCreateTime(p entql.TimeP) {
	f.Where(p.Field(tenantconn.FieldCreateTime))
}

// WhereUpdateTime applies the entql time.Time predicate on the update_time field.
func (f *TenantConnFilter) WhereUpdateTime(p entql.TimeP) {
	f.Where(p.Field(tenantconn.FieldUpdateTime))
}

// WhereKey applies the entql string predicate on the key field.
func (f *TenantConnFilter) WhereKey(p entql.StringP) {
	f.Where(p.Field(tenantconn.FieldKey))
}

// WhereValue applies the entql string predicate on the value field.
func (f *TenantConnFilter) WhereValue(p entql.StringP) {
	f.Where(p.Field(tenantconn.FieldValue))
}