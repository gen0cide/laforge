// Code generated by entc, DO NOT EDIT.

package provisionedhost

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/gen0cide/laforge/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
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
func IDGT(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// SubnetIP applies equality check predicate on the "subnet_ip" field. It's identical to SubnetIPEQ.
func SubnetIP(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPEQ applies the EQ predicate on the "subnet_ip" field.
func SubnetIPEQ(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPNEQ applies the NEQ predicate on the "subnet_ip" field.
func SubnetIPNEQ(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPIn applies the In predicate on the "subnet_ip" field.
func SubnetIPIn(vs ...string) predicate.ProvisionedHost {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSubnetIP), v...))
	})
}

// SubnetIPNotIn applies the NotIn predicate on the "subnet_ip" field.
func SubnetIPNotIn(vs ...string) predicate.ProvisionedHost {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSubnetIP), v...))
	})
}

// SubnetIPGT applies the GT predicate on the "subnet_ip" field.
func SubnetIPGT(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPGTE applies the GTE predicate on the "subnet_ip" field.
func SubnetIPGTE(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPLT applies the LT predicate on the "subnet_ip" field.
func SubnetIPLT(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPLTE applies the LTE predicate on the "subnet_ip" field.
func SubnetIPLTE(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPContains applies the Contains predicate on the "subnet_ip" field.
func SubnetIPContains(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPHasPrefix applies the HasPrefix predicate on the "subnet_ip" field.
func SubnetIPHasPrefix(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPHasSuffix applies the HasSuffix predicate on the "subnet_ip" field.
func SubnetIPHasSuffix(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPEqualFold applies the EqualFold predicate on the "subnet_ip" field.
func SubnetIPEqualFold(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldSubnetIP), v))
	})
}

// SubnetIPContainsFold applies the ContainsFold predicate on the "subnet_ip" field.
func SubnetIPContainsFold(v string) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldSubnetIP), v))
	})
}

// HasProvisionedHostToTag applies the HasEdge predicate on the "ProvisionedHostToTag" edge.
func HasProvisionedHostToTag() predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToTagTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProvisionedHostToTagTable, ProvisionedHostToTagColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionedHostToTagWith applies the HasEdge predicate on the "ProvisionedHostToTag" edge with a given conditions (other predicates).
func HasProvisionedHostToTagWith(preds ...predicate.Tag) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToTagInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProvisionedHostToTagTable, ProvisionedHostToTagColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionedHostToStatus applies the HasEdge predicate on the "ProvisionedHostToStatus" edge.
func HasProvisionedHostToStatus() predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToStatusTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProvisionedHostToStatusTable, ProvisionedHostToStatusColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionedHostToStatusWith applies the HasEdge predicate on the "ProvisionedHostToStatus" edge with a given conditions (other predicates).
func HasProvisionedHostToStatusWith(preds ...predicate.Status) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToStatusInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProvisionedHostToStatusTable, ProvisionedHostToStatusColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionedHostToProvisionedNetwork applies the HasEdge predicate on the "ProvisionedHostToProvisionedNetwork" edge.
func HasProvisionedHostToProvisionedNetwork() predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToProvisionedNetworkTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ProvisionedHostToProvisionedNetworkTable, ProvisionedHostToProvisionedNetworkPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionedHostToProvisionedNetworkWith applies the HasEdge predicate on the "ProvisionedHostToProvisionedNetwork" edge with a given conditions (other predicates).
func HasProvisionedHostToProvisionedNetworkWith(preds ...predicate.ProvisionedNetwork) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToProvisionedNetworkInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, ProvisionedHostToProvisionedNetworkTable, ProvisionedHostToProvisionedNetworkPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionedHostToHost applies the HasEdge predicate on the "ProvisionedHostToHost" edge.
func HasProvisionedHostToHost() predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToHostTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProvisionedHostToHostTable, ProvisionedHostToHostColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionedHostToHostWith applies the HasEdge predicate on the "ProvisionedHostToHost" edge with a given conditions (other predicates).
func HasProvisionedHostToHostWith(preds ...predicate.Host) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToHostInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProvisionedHostToHostTable, ProvisionedHostToHostColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionedHostToProvisioningStep applies the HasEdge predicate on the "ProvisionedHostToProvisioningStep" edge.
func HasProvisionedHostToProvisioningStep() predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToProvisioningStepTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ProvisionedHostToProvisioningStepTable, ProvisionedHostToProvisioningStepPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionedHostToProvisioningStepWith applies the HasEdge predicate on the "ProvisionedHostToProvisioningStep" edge with a given conditions (other predicates).
func HasProvisionedHostToProvisioningStepWith(preds ...predicate.ProvisioningStep) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToProvisioningStepInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ProvisionedHostToProvisioningStepTable, ProvisionedHostToProvisioningStepPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionedHostToAgentStatus applies the HasEdge predicate on the "ProvisionedHostToAgentStatus" edge.
func HasProvisionedHostToAgentStatus() predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToAgentStatusTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ProvisionedHostToAgentStatusTable, ProvisionedHostToAgentStatusPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionedHostToAgentStatusWith applies the HasEdge predicate on the "ProvisionedHostToAgentStatus" edge with a given conditions (other predicates).
func HasProvisionedHostToAgentStatusWith(preds ...predicate.AgentStatus) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(ProvisionedHostToAgentStatusInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ProvisionedHostToAgentStatusTable, ProvisionedHostToAgentStatusPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ProvisionedHost) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ProvisionedHost) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
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
func Not(p predicate.ProvisionedHost) predicate.ProvisionedHost {
	return predicate.ProvisionedHost(func(s *sql.Selector) {
		p(s.Not())
	})
}