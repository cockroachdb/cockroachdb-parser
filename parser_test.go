package main

import (
	"testing"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	p, err := parser.ParseOne("CREATE TABLE t (a TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP)")
	require.NoError(t, err)
	f := tree.DefaultPrettyCfg()
	t.Log(f.Pretty(p.AST))
}
