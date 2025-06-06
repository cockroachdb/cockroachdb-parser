// Copyright 2018 The Cockroach Authors.
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

const vecTmpl = "pkg/col/coldata/vec_tmpl.go"

func genVec(inputFileContents string, wr io.Writer) error {
	r := strings.NewReplacer("_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
		"_TYPE_WIDTH", typeWidthReplacement,
		"_GOTYPESLICE", "{{.GoTypeSliceNameInColdata}}",
		"_GOTYPE", "{{.GoType}}",
		"_TYPE", "{{.VecMethod}}",
		"TemplateType", "{{.VecMethod}}",
	)
	s := r.Replace(inputFileContents)

	copyWithReorderedSource := makeFunctionRegex("_COPY_WITH_REORDERED_SOURCE", 1)
	s = copyWithReorderedSource.ReplaceAllString(s, `{{template "copyWithReorderedSource" buildDict "Global" . "SrcHasNulls" $1}}`)

	s = replaceManipulationFuncs(s)

	// Now, generate the op, from the template.
	tmpl, err := template.New("vec_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}

	// It doesn't matter that we're passing in all overloads of Equality
	// comparison operator - we simply need to iterate over all supported
	// types.
	return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.EQ])
}
func init() {
	registerGenerator(genVec, "vec.eg.go", vecTmpl)
}
