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

const minMaxRemovableAggTmpl = "pkg/sql/colexec/colexecwindow/min_max_removable_agg_tmpl.go"

func genMinMaxRemovableAgg(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer(
		"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_AGG_TITLE", "{{.AggTitle}}",
		"_AGG", "{{$agg}}",
		"_GOTYPE", "{{.GoType}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	assignCmpRe := makeFunctionRegex("_ASSIGN_CMP", 6)
	s = assignCmpRe.ReplaceAllString(s, makeTemplateFunctionCall("Assign", 6))

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("min_max_agg").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, []struct {
		Agg       string
		AggTitle  string
		Overloads []*oneArgOverload
	}{
		{
			Agg:       "min",
			AggTitle:  "Min",
			Overloads: sameTypeComparisonOpToOverloads[treecmp.LT],
		},
		{
			Agg:       "max",
			AggTitle:  "Max",
			Overloads: sameTypeComparisonOpToOverloads[treecmp.GT],
		},
	})
}

func init() {
	registerGenerator(genMinMaxRemovableAgg, "min_max_removable_agg.eg.go", minMaxRemovableAggTmpl)
}
