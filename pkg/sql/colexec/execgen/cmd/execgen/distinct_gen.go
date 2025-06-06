// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree/treecmp"
)

// distinctTmpl is the common base for the template used to generate code for
// ordered distinct structs and sort partitioners. It should be used as a format
// string with two %s arguments:
// 1. specifies the package that the generated code is placed in
// 2. specifies which of the building blocks coming from distinct_tmpl.go should
// be included into the code generation.
const distinctTmpl = `
// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package %s

import (
	"context"

	"github.com/cockroachdb/apd/v3"
	"github.com/cockroachdb/cockroachdb-parser/pkg/col/coldata"
	"github.com/cockroachdb/cockroachdb-parser/pkg/col/coldataext"
	"github.com/cockroachdb/cockroachdb-parser/pkg/col/typeconv"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/colexec/execgen"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/colexecop"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/colexecerror"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/duration"
	"github.com/cockroachdb/cockroachdb-parser/pkg/util/json"
	"github.com/cockroachdb/errors"
)

%s
`

const distinctOpsTmpl = "pkg/sql/colexec/colexecbase/distinct_tmpl.go"

func genDistinctOps(targetPkg, targetTmpl string) generator {
	return func(inputFileContents string, wr io.Writer) error {
		r := strings.NewReplacer(
			"_CANONICAL_TYPE_FAMILY", "{{.CanonicalTypeFamilyStr}}",
			"_TYPE_WIDTH", typeWidthReplacement,
			"_GOTYPESLICE", "{{.GoTypeSliceName}}",
			"_GOTYPE", "{{.GoType}}",
			"_TYPE", "{{.VecMethod}}",
			"TemplateType", "{{.VecMethod}}")
		s := r.Replace(inputFileContents)

		assignNeRe := makeFunctionRegex("_ASSIGN_NE", 6)
		s = assignNeRe.ReplaceAllString(s, makeTemplateFunctionCall("Assign", 6))

		s = replaceManipulationFuncs(s)

		// Now, generate the op, from the template.
		tmpl, err := template.New("distinct_op").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
		if err != nil {
			return err
		}

		tmpl, err = tmpl.Parse(fmt.Sprintf(distinctTmpl, targetPkg, targetTmpl))
		if err != nil {
			return err
		}
		return tmpl.Execute(wr, sameTypeComparisonOpToOverloads[treecmp.NE])
	}
}

func init() {
	distinctOp := `
{{template "distinctOpConstructor" .}}

{{range .}}
{{range .WidthOverloads}}
{{template "distinctOp" .}}
{{end}}
{{end}}
`
	sortPartitioner := `
{{template "sortPartitionerConstructor" .}}

{{range .}}
{{range .WidthOverloads}}
{{template "sortPartitioner" .}}
{{end}}
{{end}}
`
	registerGenerator(genDistinctOps("colexecbase", distinctOp), "distinct.eg.go", distinctOpsTmpl)
	registerGenerator(genDistinctOps("colexec", sortPartitioner), "sort_partitioner.eg.go", distinctOpsTmpl)
}
