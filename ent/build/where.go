// Code generated by entc, DO NOT EDIT.

package build

import (
	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/gen0cide/laforge/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Revision applies equality check predicate on the "revision" field. It's identical to RevisionEQ.
func Revision(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRevision), v))
	})
}

// RevisionEQ applies the EQ predicate on the "revision" field.
func RevisionEQ(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRevision), v))
	})
}

// RevisionNEQ applies the NEQ predicate on the "revision" field.
func RevisionNEQ(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRevision), v))
	})
}

// RevisionIn applies the In predicate on the "revision" field.
func RevisionIn(vs ...int) predicate.Build {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Build(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldRevision), v...))
	})
}

// RevisionNotIn applies the NotIn predicate on the "revision" field.
func RevisionNotIn(vs ...int) predicate.Build {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Build(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldRevision), v...))
	})
}

// RevisionGT applies the GT predicate on the "revision" field.
func RevisionGT(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRevision), v))
	})
}

// RevisionGTE applies the GTE predicate on the "revision" field.
func RevisionGTE(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRevision), v))
	})
}

// RevisionLT applies the LT predicate on the "revision" field.
func RevisionLT(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRevision), v))
	})
}

// RevisionLTE applies the LTE predicate on the "revision" field.
func RevisionLTE(v int) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRevision), v))
	})
}

// HasBuildToUser applies the HasEdge predicate on the "BuildToUser" edge.
func HasBuildToUser() predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToUserTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BuildToUserTable, BuildToUserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBuildToUserWith applies the HasEdge predicate on the "BuildToUser" edge with a given conditions (other predicates).
func HasBuildToUserWith(preds ...predicate.User) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToUserInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BuildToUserTable, BuildToUserColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBuildToTag applies the HasEdge predicate on the "BuildToTag" edge.
func HasBuildToTag() predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToTagTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BuildToTagTable, BuildToTagColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBuildToTagWith applies the HasEdge predicate on the "BuildToTag" edge with a given conditions (other predicates).
func HasBuildToTagWith(preds ...predicate.Tag) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToTagInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, BuildToTagTable, BuildToTagColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBuildToProvisionedNetwork applies the HasEdge predicate on the "BuildToProvisionedNetwork" edge.
func HasBuildToProvisionedNetwork() predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToProvisionedNetworkTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, BuildToProvisionedNetworkTable, BuildToProvisionedNetworkPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBuildToProvisionedNetworkWith applies the HasEdge predicate on the "BuildToProvisionedNetwork" edge with a given conditions (other predicates).
func HasBuildToProvisionedNetworkWith(preds ...predicate.ProvisionedNetwork) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToProvisionedNetworkInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, BuildToProvisionedNetworkTable, BuildToProvisionedNetworkPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBuildToTeam applies the HasEdge predicate on the "BuildToTeam" edge.
func HasBuildToTeam() predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToTeamTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, BuildToTeamTable, BuildToTeamPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBuildToTeamWith applies the HasEdge predicate on the "BuildToTeam" edge with a given conditions (other predicates).
func HasBuildToTeamWith(preds ...predicate.Team) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToTeamInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, BuildToTeamTable, BuildToTeamPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasBuildToEnvironment applies the HasEdge predicate on the "BuildToEnvironment" edge.
func HasBuildToEnvironment() predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToEnvironmentTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, BuildToEnvironmentTable, BuildToEnvironmentPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasBuildToEnvironmentWith applies the HasEdge predicate on the "BuildToEnvironment" edge with a given conditions (other predicates).
func HasBuildToEnvironmentWith(preds ...predicate.Environment) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(BuildToEnvironmentInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, BuildToEnvironmentTable, BuildToEnvironmentPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Build) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Build) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Build) predicate.Build {
	return predicate.Build(func(s *sql.Selector) {
		p(s.Not())
	})
}
