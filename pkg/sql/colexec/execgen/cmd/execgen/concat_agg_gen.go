// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package main

import (
	"io"
	"text/template"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"
)

const concatAggTmpl = "pkg/sql/colexec/colexecagg/concat_agg_tmpl.go"

func genConcatAgg(inputFileContents string, wr io.Writer) error {
	accumulateConcatRe := makeFunctionRegex("_ACCUMULATE_CONCAT", 5)
	s := accumulateConcatRe.ReplaceAllString(inputFileContents, `{{template "accumulateConcat" buildDict "HasNulls" $4 "HasSel" $5}}`)

	s = replaceManipulationFuncs(s)

	tmpl, err := template.New("concat_agg").Funcs(template.FuncMap{"buildDict": buildDict}).Parse(s)
	if err != nil {
		return err
	}
	return tmpl.Execute(wr, aggTmplInfoBase{canonicalTypeFamily: types.BytesFamily})
}

func init() {
	registerAggGenerator(
		genConcatAgg, "concat_agg.eg.go", /* filenameSuffix */
		concatAggTmpl, "concat" /* aggName */, true, /* genWindowVariant */
	)
}
