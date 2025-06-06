// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package execgen

import (
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree/treebin"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree/treecmp"
)

// BinaryOpName is a mapping from all binary operators that are supported by
// the vectorized engine to their names.
var BinaryOpName = map[treebin.BinaryOperatorSymbol]string{
	treebin.Bitand:            "Bitand",
	treebin.Bitor:             "Bitor",
	treebin.Bitxor:            "Bitxor",
	treebin.Plus:              "Plus",
	treebin.Minus:             "Minus",
	treebin.Mult:              "Mult",
	treebin.Div:               "Div",
	treebin.FloorDiv:          "FloorDiv",
	treebin.Mod:               "Mod",
	treebin.Pow:               "Pow",
	treebin.Concat:            "Concat",
	treebin.LShift:            "LShift",
	treebin.RShift:            "RShift",
	treebin.JSONFetchVal:      "JSONFetchVal",
	treebin.JSONFetchText:     "JSONFetchText",
	treebin.JSONFetchValPath:  "JSONFetchValPath",
	treebin.JSONFetchTextPath: "JSONFetchTextPath",
}

// ComparisonOpName is a mapping from all comparison operators that are
// supported by the vectorized engine to their names.
var ComparisonOpName = map[treecmp.ComparisonOperatorSymbol]string{
	treecmp.EQ: "EQ",
	treecmp.NE: "NE",
	treecmp.LT: "LT",
	treecmp.LE: "LE",
	treecmp.GT: "GT",
	treecmp.GE: "GE",
}
