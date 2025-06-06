// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package main

import (
	"io"
	"strings"
	"text/template"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree/treecmp"
)

const crossJoinerTmpl = "pkg/sql/colexec/colexecjoin/crossjoiner_tmpl.go"

func genCrossJoiner(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("crossjoiner").Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}

func init() {
	registerGenerator(genCrossJoiner, "crossjoiner.eg.go", crossJoinerTmpl)
}
