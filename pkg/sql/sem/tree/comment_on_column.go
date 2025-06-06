// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package tree

import "github.com/cockroachdb/cockroachdb-parser/pkg/sql/lexbase"

// CommentOnColumn represents an COMMENT ON COLUMN statement.
type CommentOnColumn struct {
	*ColumnItem
	Comment *string
}

// Format implements the NodeFormatter interface.
func (n *CommentOnColumn) Format(ctx *FmtCtx) {
	ctx.WriteString("COMMENT ON COLUMN ")
	ctx.FormatNode(n.ColumnItem)
	ctx.WriteString(" IS ")
	if n.Comment != nil {
		// TODO(knz): Replace all this with ctx.FormatNode
		// when COMMENT supports expressions.
		if ctx.flags.HasFlags(FmtHideConstants) {
			ctx.WriteString("'_'")
		} else {
			lexbase.EncodeSQLStringWithFlags(&ctx.Buffer, *n.Comment, ctx.flags.EncodeFlags())
		}
	} else {
		ctx.WriteString("NULL")
	}
}
