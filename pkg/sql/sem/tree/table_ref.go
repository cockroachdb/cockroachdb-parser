// Copyright 2017 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package tree

import "github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/catid"

// ID is a custom type for {Database,Table}Descriptor IDs.
type ID = catid.ColumnID

// ColumnID is a custom type for ColumnDescriptor IDs.
type ColumnID = catid.ColumnID

// TableRef represents a numeric table reference.
// (Syntax !NNN in SQL.)
type TableRef struct {
	// TableID is the descriptor ID of the requested table.
	TableID int64

	// ColumnIDs is the list of column IDs requested in the table.
	// Note that a nil array here means "unspecified" (all columns)
	// whereas an array of length 0 means "zero columns".
	// Lists of zero columns are not supported and will throw an error.
	Columns []ColumnID

	// As determines the names that can be used in the surrounding query
	// to refer to this source.
	As AliasClause
}

// Format implements the NodeFormatter interface.
func (n *TableRef) Format(ctx *FmtCtx) {
	ctx.Printf("[%d", n.TableID)
	if n.Columns != nil {
		ctx.WriteByte('(')
		for i, c := range n.Columns {
			if i > 0 {
				ctx.WriteString(", ")
			}
			ctx.Printf("%d", c)
		}
		ctx.WriteByte(')')
	}
	if n.As.Alias != "" {
		ctx.WriteString(" AS ")
		ctx.FormatNode(&n.As)
	}
	ctx.WriteByte(']')
}
func (n *TableRef) String() string { return AsString(n) }

// tableExpr implements the TableExpr interface.
func (n *TableRef) tableExpr() {}
