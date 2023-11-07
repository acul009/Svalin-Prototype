// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"rahnit-rmm/ent/predicate"
	"rahnit-rmm/ent/tunnelconfig"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// TunnelConfigDelete is the builder for deleting a TunnelConfig entity.
type TunnelConfigDelete struct {
	config
	hooks    []Hook
	mutation *TunnelConfigMutation
}

// Where appends a list predicates to the TunnelConfigDelete builder.
func (tcd *TunnelConfigDelete) Where(ps ...predicate.TunnelConfig) *TunnelConfigDelete {
	tcd.mutation.Where(ps...)
	return tcd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tcd *TunnelConfigDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tcd.sqlExec, tcd.mutation, tcd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tcd *TunnelConfigDelete) ExecX(ctx context.Context) int {
	n, err := tcd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tcd *TunnelConfigDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(tunnelconfig.Table, sqlgraph.NewFieldSpec(tunnelconfig.FieldID, field.TypeInt))
	if ps := tcd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tcd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tcd.mutation.done = true
	return affected, err
}

// TunnelConfigDeleteOne is the builder for deleting a single TunnelConfig entity.
type TunnelConfigDeleteOne struct {
	tcd *TunnelConfigDelete
}

// Where appends a list predicates to the TunnelConfigDelete builder.
func (tcdo *TunnelConfigDeleteOne) Where(ps ...predicate.TunnelConfig) *TunnelConfigDeleteOne {
	tcdo.tcd.mutation.Where(ps...)
	return tcdo
}

// Exec executes the deletion query.
func (tcdo *TunnelConfigDeleteOne) Exec(ctx context.Context) error {
	n, err := tcdo.tcd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{tunnelconfig.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tcdo *TunnelConfigDeleteOne) ExecX(ctx context.Context) {
	if err := tcdo.Exec(ctx); err != nil {
		panic(err)
	}
}