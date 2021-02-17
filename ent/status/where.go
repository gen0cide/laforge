// Code generated by entc, DO NOT EDIT.

package status

import (
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/gen0cide/laforge/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
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
func IDGT(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// StartedAt applies equality check predicate on the "started_at" field. It's identical to StartedAtEQ.
func StartedAt(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStartedAt), v))
	})
}

// EndedAt applies equality check predicate on the "ended_at" field. It's identical to EndedAtEQ.
func EndedAt(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEndedAt), v))
	})
}

// Failed applies equality check predicate on the "failed" field. It's identical to FailedEQ.
func Failed(v bool) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFailed), v))
	})
}

// Completed applies equality check predicate on the "completed" field. It's identical to CompletedEQ.
func Completed(v bool) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCompleted), v))
	})
}

// Error applies equality check predicate on the "error" field. It's identical to ErrorEQ.
func Error(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldError), v))
	})
}

// StateEQ applies the EQ predicate on the "state" field.
func StateEQ(v State) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldState), v))
	})
}

// StateNEQ applies the NEQ predicate on the "state" field.
func StateNEQ(v State) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldState), v))
	})
}

// StateIn applies the In predicate on the "state" field.
func StateIn(vs ...State) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldState), v...))
	})
}

// StateNotIn applies the NotIn predicate on the "state" field.
func StateNotIn(vs ...State) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldState), v...))
	})
}

// StartedAtEQ applies the EQ predicate on the "started_at" field.
func StartedAtEQ(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStartedAt), v))
	})
}

// StartedAtNEQ applies the NEQ predicate on the "started_at" field.
func StartedAtNEQ(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStartedAt), v))
	})
}

// StartedAtIn applies the In predicate on the "started_at" field.
func StartedAtIn(vs ...time.Time) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStartedAt), v...))
	})
}

// StartedAtNotIn applies the NotIn predicate on the "started_at" field.
func StartedAtNotIn(vs ...time.Time) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStartedAt), v...))
	})
}

// StartedAtGT applies the GT predicate on the "started_at" field.
func StartedAtGT(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStartedAt), v))
	})
}

// StartedAtGTE applies the GTE predicate on the "started_at" field.
func StartedAtGTE(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStartedAt), v))
	})
}

// StartedAtLT applies the LT predicate on the "started_at" field.
func StartedAtLT(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStartedAt), v))
	})
}

// StartedAtLTE applies the LTE predicate on the "started_at" field.
func StartedAtLTE(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStartedAt), v))
	})
}

// EndedAtEQ applies the EQ predicate on the "ended_at" field.
func EndedAtEQ(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldEndedAt), v))
	})
}

// EndedAtNEQ applies the NEQ predicate on the "ended_at" field.
func EndedAtNEQ(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldEndedAt), v))
	})
}

// EndedAtIn applies the In predicate on the "ended_at" field.
func EndedAtIn(vs ...time.Time) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldEndedAt), v...))
	})
}

// EndedAtNotIn applies the NotIn predicate on the "ended_at" field.
func EndedAtNotIn(vs ...time.Time) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldEndedAt), v...))
	})
}

// EndedAtGT applies the GT predicate on the "ended_at" field.
func EndedAtGT(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldEndedAt), v))
	})
}

// EndedAtGTE applies the GTE predicate on the "ended_at" field.
func EndedAtGTE(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldEndedAt), v))
	})
}

// EndedAtLT applies the LT predicate on the "ended_at" field.
func EndedAtLT(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldEndedAt), v))
	})
}

// EndedAtLTE applies the LTE predicate on the "ended_at" field.
func EndedAtLTE(v time.Time) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldEndedAt), v))
	})
}

// FailedEQ applies the EQ predicate on the "failed" field.
func FailedEQ(v bool) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldFailed), v))
	})
}

// FailedNEQ applies the NEQ predicate on the "failed" field.
func FailedNEQ(v bool) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldFailed), v))
	})
}

// CompletedEQ applies the EQ predicate on the "completed" field.
func CompletedEQ(v bool) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCompleted), v))
	})
}

// CompletedNEQ applies the NEQ predicate on the "completed" field.
func CompletedNEQ(v bool) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCompleted), v))
	})
}

// ErrorEQ applies the EQ predicate on the "error" field.
func ErrorEQ(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldError), v))
	})
}

// ErrorNEQ applies the NEQ predicate on the "error" field.
func ErrorNEQ(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldError), v))
	})
}

// ErrorIn applies the In predicate on the "error" field.
func ErrorIn(vs ...string) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldError), v...))
	})
}

// ErrorNotIn applies the NotIn predicate on the "error" field.
func ErrorNotIn(vs ...string) predicate.Status {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Status(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldError), v...))
	})
}

// ErrorGT applies the GT predicate on the "error" field.
func ErrorGT(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldError), v))
	})
}

// ErrorGTE applies the GTE predicate on the "error" field.
func ErrorGTE(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldError), v))
	})
}

// ErrorLT applies the LT predicate on the "error" field.
func ErrorLT(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldError), v))
	})
}

// ErrorLTE applies the LTE predicate on the "error" field.
func ErrorLTE(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldError), v))
	})
}

// ErrorContains applies the Contains predicate on the "error" field.
func ErrorContains(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldError), v))
	})
}

// ErrorHasPrefix applies the HasPrefix predicate on the "error" field.
func ErrorHasPrefix(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldError), v))
	})
}

// ErrorHasSuffix applies the HasSuffix predicate on the "error" field.
func ErrorHasSuffix(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldError), v))
	})
}

// ErrorEqualFold applies the EqualFold predicate on the "error" field.
func ErrorEqualFold(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldError), v))
	})
}

// ErrorContainsFold applies the ContainsFold predicate on the "error" field.
func ErrorContainsFold(v string) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldError), v))
	})
}

// HasStatusToTag applies the HasEdge predicate on the "StatusToTag" edge.
func HasStatusToTag() predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StatusToTagTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, StatusToTagTable, StatusToTagColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStatusToTagWith applies the HasEdge predicate on the "StatusToTag" edge with a given conditions (other predicates).
func HasStatusToTagWith(preds ...predicate.Tag) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StatusToTagInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, StatusToTagTable, StatusToTagColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Status) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Status) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
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
func Not(p predicate.Status) predicate.Status {
	return predicate.Status(func(s *sql.Selector) {
		p(s.Not())
	})
}
