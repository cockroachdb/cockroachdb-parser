package main

import (
	"testing"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	_, err := parser.Parse("SELECT 1")
	require.NoError(t, err)
}
