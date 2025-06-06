// Copyright 2021 The Cockroach Authors.
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

const leadLagTmpl = "pkg/sql/colexec/colexecwindow/lead_lag_tmpl.go"

func replaceLeadLagTmplVariables(tmpl string) string {
	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_GOTYPESLICE", "{{.GoTypeSliceName}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(tmpl)

	processBatch := makeFunctionRegex("_PROCESS_BATCH", 2)
	s = processBatch.ReplaceAllString(s, `{{template "processBatchTmpl" buildDict "VecMethod" .VecMethod "IsBytesLike" .IsBytesLike "OffsetHasNulls" $1 "DefaultHasNulls" $2}}`)

	return replaceManipulationFuncs(s)
}

func genLeadOp(inputFileContents string, wr io.Writer) error {
	s := replaceLeadLagTmplVariables(inputFileContents)
	r := strings.NewReplacer(
		"_OP_NAME", "lead",
		"_UPPERCASE_NAME", "Lead",
	)
	s = r.Replace(s)

	// Now, generate the op, from the template.
	tmpl, err := template.New("lead_ops").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}

func genLagOp(inputFileContents string, wr io.Writer) error {
	s := replaceLeadLagTmplVariables(inputFileContents)
	r := strings.NewReplacer(
		"_OP_NAME", "lag",
		"_UPPERCASE_NAME", "Lag",
	)
	s = r.Replace(s)

	// Now, generate the op, from the template.
	tmpl, err := template.New("lag_ops").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}

func init() {
	registerGenerator(genLeadOp, "lead.eg.go", leadLagTmpl)
	registerGenerator(genLagOp, "lag.eg.go", leadLagTmpl)
}
