// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in licenses/BSD-vitess.txt.

// Portions of this file are additionally subject to the following
// license and copyright.
//
// Copyright 2015 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

// This code was derived from https://github.com/youtube/vitess.

package parser

import (
	"fmt"
	"go/constant"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser/statements"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/scanner"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/types"
	"github.com/cockroachdb/errors"
)

func init() {
	scanner.NewNumValFn = func(a constant.Value, s string, b bool) interface{} { return tree.NewNumVal(a, s, b) }
	scanner.NewPlaceholderFn = func(s string) (interface{}, error) { return tree.NewPlaceholder(s) }
}

// Parser wraps a scanner, parser and other utilities present in the parser
// package.
type Parser struct {
	scanner    scanner.SQLScanner
	lexer      lexer
	parserImpl sqlParserImpl
	tokBuf     [8]sqlSymType
	stmtBuf    [1]statements.Statement[tree.Statement]
}

// INT8 is the historical interpretation of INT. This should be left
// alone in the future, since there are many sql fragments stored
// in various descriptors. Any user input that was created after
// INT := INT4 will simply use INT4 in any resulting code.
var defaultNakedIntType = types.Int

// NakedIntTypeFromDefaultIntSize given the size in bits or bytes (preferred)
// of how a "naked" INT type should be parsed returns the corresponding integer
// type.
func NakedIntTypeFromDefaultIntSize(defaultIntSize int32) *types.T {
	switch defaultIntSize {
	case 4, 32:
		return types.Int4
	default:
		return types.Int
	}
}

// Parse parses the sql and returns a list of statements.
func (p *Parser) Parse(sql string) (statements.Statements, error) {
	return p.parseWithDepth(1, sql, defaultNakedIntType, discardComments)
}

// ParseWithInt parses a sql statement string and returns a list of
// Statements. The INT token will result in the specified TInt type.
func (p *Parser) ParseWithInt(sql string, nakedIntType *types.T) (statements.Statements, error) {
	return p.parseWithDepth(1, sql, nakedIntType, discardComments)
}

func (p *Parser) parseOneWithInt(
	sql string, nakedIntType *types.T, comments commentsMode,
) (statements.Statement[tree.Statement], error) {
	stmts, err := p.parseWithDepth(1, sql, nakedIntType, comments)
	if err != nil {
		return statements.Statement[tree.Statement]{}, err
	}
	if len(stmts) != 1 {
		return statements.Statement[tree.Statement]{}, errors.AssertionFailedf("expected 1 statement, but found %d", len(stmts))
	}
	return stmts[0], nil
}

func (p *Parser) scanOneStmt() (sql string, tokens []sqlSymType, done bool) {
	tokens = p.tokBuf[:0]
	tokens = append(tokens, sqlSymType{})
	lval := &p.tokBuf[0]

	// Scan the first token.
	for {
		p.scanner.Scan(lval)
		if lval.id == 0 {
			return "", nil, true
		}
		if lval.id != ';' {
			break
		}
	}

	startPos := lval.pos
	// We make the resulting token positions match the returned string.
	lval.pos = 0
	var preValID int32
	// This is used to track the degree of nested `BEGIN ATOMIC ... END` function
	// body context. When greater than zero, it means that we're scanning through
	// the function body of a `CREATE FUNCTION` statement. ';' character is only
	// a separator of sql statements within the body instead of a finishing line
	// of the `CREATE FUNCTION` statement.
	curFuncBodyCnt := 0
	for {
		if lval.id == ERROR {
			return p.scanner.In()[startPos:], tokens, true
		}
		preValID = lval.id
		tokens = append(tokens, sqlSymType{})
		lval = &tokens[len(tokens)-1]
		p.scanner.Scan(lval)

		if preValID == BEGIN && lval.id == ATOMIC {
			curFuncBodyCnt++
		}
		if curFuncBodyCnt > 0 && lval.id == END {
			curFuncBodyCnt--
		}
		if lval.id == 0 || (curFuncBodyCnt == 0 && lval.id == ';') {
			endPos := p.scanner.Pos()
			if lval.id == ';' {
				// Don't include the ending semicolon, if there is one, in the raw SQL.
				endPos--
			}
			tokens = tokens[:len(tokens)-1]
			return p.scanner.In()[startPos:endPos], tokens, (lval.id == 0)
		}
		lval.pos -= startPos
	}
}

type commentsMode bool

const (
	retainComments  commentsMode = true
	discardComments commentsMode = false
)

func (p *Parser) parseWithDepth(
	depth int, sql string, nakedIntType *types.T, cm commentsMode,
) (statements.Statements, error) {
	stmts := statements.Statements(p.stmtBuf[:0])
	p.scanner.Init(sql)
	if cm == retainComments {
		p.scanner.RetainComments()
	}
	defer p.scanner.Cleanup()
	for {
		sql, tokens, done := p.scanOneStmt()
		stmt, err := p.parse(depth+1, sql, tokens, nakedIntType)
		if err != nil {
			return nil, err
		}
		if stmt.AST != nil {
			stmts = append(stmts, stmt)
		}
		if done {
			break
		}
	}
	return stmts, nil
}

// parse parses a statement from the given scanned tokens.
func (p *Parser) parse(
	depth int, sql string, tokens []sqlSymType, nakedIntType *types.T,
) (statements.Statement[tree.Statement], error) {
	p.lexer.init(sql, tokens, nakedIntType)
	defer p.lexer.cleanup()
	if p.parserImpl.Parse(&p.lexer) != 0 {
		if p.lexer.lastError == nil {
			// This should never happen -- there should be an error object
			// every time Parse() returns nonzero. We're just playing safe
			// here.
			p.lexer.Error("syntax error")
		}
		err := p.lexer.lastError

		// Compatibility with 19.1 telemetry: prefix the telemetry keys
		// with the "syntax." prefix.
		// TODO(knz): move the auto-prefixing of feature names to a
		// higher level in the call stack.
		tkeys := errors.GetTelemetryKeys(err)
		if len(tkeys) > 0 {
			for i := range tkeys {
				tkeys[i] = "syntax." + tkeys[i]
			}
			err = errors.WithTelemetry(err, tkeys...)
		}

		return statements.Statement[tree.Statement]{}, err
	}

	return statements.Statement[tree.Statement]{
		AST:             p.lexer.stmt,
		SQL:             sql,
		Comments:        p.scanner.Comments,
		NumPlaceholders: p.lexer.numPlaceholders,
		NumAnnotations:  p.lexer.numAnnotations,
	}, nil
}

// unaryNegation constructs an AST node for a negation. This attempts
// to preserve constant NumVals and embed the negative sign inside
// them instead of wrapping in an UnaryExpr. This in turn ensures
// that negative numbers get considered as a single constant
// for the purpose of formatting and scrubbing.
func unaryNegation(e tree.Expr) tree.Expr {
	if cst, ok := e.(*tree.NumVal); ok {
		cst.Negate()
		return cst
	}

	// Common case.
	return &tree.UnaryExpr{
		Operator: tree.MakeUnaryOperator(tree.UnaryMinus),
		Expr:     e,
	}
}

// Parse parses a sql statement string and returns a list of Statements.
func Parse(sql string) (statements.Statements, error) {
	return ParseWithInt(sql, defaultNakedIntType)
}

// ParseWithInt parses a sql statement string and returns a list of
// Statements. The INT token will result in the specified TInt type.
func ParseWithInt(sql string, nakedIntType *types.T) (statements.Statements, error) {
	var p Parser
	return p.parseWithDepth(1, sql, nakedIntType, discardComments)
}

// ParseOne parses a sql statement string, ensuring that it contains only a
// single statement, and returns that Statement. ParseOne will always
// interpret the INT and SERIAL types as 64-bit types, since this is
// used in various internal-execution paths where we might receive
// bits of SQL from other nodes. In general, we expect that all
// user-generated SQL has been run through the ParseWithInt() function.
func ParseOne(sql string) (statements.Statement[tree.Statement], error) {
	return ParseOneWithInt(sql, defaultNakedIntType)
}

// ParseOneRetainComments is similar to ParseOne, but it retains scanned
// comments in the returned statement's Comment field.
func ParseOneRetainComments(sql string) (statements.Statement[tree.Statement], error) {
	var p Parser
	return p.parseOneWithInt(sql, defaultNakedIntType, retainComments)
}

// ParseOneWithInt is similar to ParseOn but interprets the INT and SERIAL
// types as the provided integer type.
func ParseOneWithInt(
	sql string, nakedIntType *types.T,
) (statements.Statement[tree.Statement], error) {
	var p Parser
	return p.parseOneWithInt(sql, nakedIntType, discardComments)
}

// ParseQualifiedTableName parses a possibly qualified table name. The
// table name must contain one or more name parts, using the full
// input SQL syntax: each name part containing special characters, or
// non-lowercase characters, must be enclosed in double quote. The
// name may not be an invalid table name (the caller is responsible
// for guaranteeing that only valid table names are provided as
// input).
func ParseQualifiedTableName(sql string) (*tree.TableName, error) {
	name, err := ParseTableName(sql)
	if err != nil {
		return nil, err
	}
	tn := name.ToTableName()
	return &tn, nil
}

// ParseTableName parses a table name. The table name must contain one
// or more name parts, using the full input SQL syntax: each name
// part containing special characters, or non-lowercase characters,
// must be enclosed in double quote. The name may not be an invalid
// table name (the caller is responsible for guaranteeing that only
// valid table names are provided as input).
func ParseTableName(sql string) (*tree.UnresolvedObjectName, error) {
	// We wrap the name we want to parse into a dummy statement since our parser
	// can only parse full statements.
	stmt, err := ParseOne(fmt.Sprintf("ALTER TABLE %s RENAME TO x", sql))
	if err != nil {
		return nil, err
	}
	rename, ok := stmt.AST.(*tree.RenameTable)
	if !ok {
		return nil, errors.AssertionFailedf("expected an ALTER TABLE statement, but found %T", stmt)
	}
	return rename.Name, nil
}

// ParseFunctionName parses a function name. The function name must contain one
// or more name parts, using the full input SQL syntax: each name
// part containing special characters, or non-lowercase characters,
// must be enclosed in double quote. The name may not be an invalid
// function name (the caller is responsible for guaranteeing that only
// valid function names are provided as input).
func ParseFunctionName(sql string) (*tree.UnresolvedObjectName, error) {
	// We wrap the name we want to parse into a dummy statement since our parser
	// can only parse full statements.
	stmt, err := ParseOne(fmt.Sprintf("ALTER FUNCTION %s RENAME TO x", sql))
	if err != nil {
		return nil, err
	}
	rename, ok := stmt.AST.(*tree.AlterRoutineRename)
	if !ok {
		return nil, errors.AssertionFailedf("expected an ALTER FUNCTION statement, but found %T", stmt)
	}
	return rename.Function.FuncName.ToUnresolvedObjectName(), nil
}

// ParseTablePattern parses a table pattern. The table name must contain one
// or more name parts, using the full input SQL syntax: each name
// part containing special characters, or non-lowercase characters,
// must be enclosed in double quote. The name may not be an invalid
// table name (the caller is responsible for guaranteeing that only
// valid table names are provided as input).
// The last part may be '*' to denote a wildcard.
func ParseTablePattern(sql string) (tree.TablePattern, error) {
	// We wrap the name we want to parse into a dummy statement since our parser
	// can only parse full statements.
	stmt, err := ParseOne(fmt.Sprintf("GRANT SELECT ON TABLE %s TO admin", sql))
	if err != nil {
		return nil, err
	}
	grant, ok := stmt.AST.(*tree.Grant)
	if !ok {
		return nil, errors.AssertionFailedf("expected a GRANT statement, but found %T", stmt)
	}
	if len(grant.Targets.Tables.TablePatterns) == 0 {
		return nil, errors.AssertionFailedf("expected at least one pattern")
	}
	u := grant.Targets.Tables.TablePatterns[0]
	un, ok := u.(*tree.UnresolvedName)
	if !ok {
		return nil, errors.AssertionFailedf("expected an unresolved name, but found %T", u)
	}
	return un.NormalizeTablePattern()
}

// parseExprsWithInt parses one or more sql expressions.
func parseExprsWithInt(exprs []string, nakedIntType *types.T) (tree.Exprs, error) {
	stmt, err := ParseOneWithInt(fmt.Sprintf("SET ROW (%s)", strings.Join(exprs, ",")), nakedIntType)
	if err != nil {
		return nil, err
	}
	set, ok := stmt.AST.(*tree.SetVar)
	if !ok {
		return nil, errors.AssertionFailedf("expected a SET statement, but found %T", stmt)
	}
	return set.Values, nil
}

// ParseExprs parses a comma-delimited sequence of SQL scalar
// expressions. The caller is responsible for ensuring that the input
// is, in fact, a comma-delimited sequence of SQL scalar expressions —
// the results are undefined if the string contains invalid SQL
// syntax.
func ParseExprs(sql []string) (tree.Exprs, error) {
	if len(sql) == 0 {
		return tree.Exprs{}, nil
	}
	return parseExprsWithInt(sql, defaultNakedIntType)
}

// ParseExpr parses a SQL scalar expression. The caller is responsible
// for ensuring that the input is, in fact, a valid SQL scalar
// expression — the results are undefined if the string contains
// invalid SQL syntax.
func ParseExpr(sql string) (tree.Expr, error) {
	return ParseExprWithInt(sql, defaultNakedIntType)
}

// ParseExprWithInt parses a SQL scalar expression, using the given
// type when INT is used as type name in the SQL syntax. The caller is
// responsible for ensuring that the input is, in fact, a valid SQL
// scalar expression — the results are undefined if the string
// contains invalid SQL syntax.
func ParseExprWithInt(sql string, nakedIntType *types.T) (tree.Expr, error) {
	exprs, err := parseExprsWithInt([]string{sql}, nakedIntType)
	if err != nil {
		return nil, err
	}
	if len(exprs) != 1 {
		return nil, errors.AssertionFailedf("expected 1 expression, found %d", len(exprs))
	}
	return exprs[0], nil
}

// GetTypeReferenceFromName turns a type name into a type
// reference. This supports only “simple” (single-identifier)
// references to built-in types, when the identifer has already been
// parsed away from the input SQL syntax.
func GetTypeReferenceFromName(typeName tree.Name) (tree.ResolvableTypeReference, error) {
	expr, err := ParseExpr(fmt.Sprintf("1::%s", typeName.String()))
	if err != nil {
		return nil, err
	}

	cast, ok := expr.(*tree.CastExpr)
	if !ok {
		return nil, errors.AssertionFailedf("expected a tree.CastExpr, but found %T", expr)
	}

	return cast.Type, nil
}

// GetTypeFromValidSQLSyntax retrieves a type from its SQL syntax. The caller is
// responsible for guaranteeing that the type expression is valid
// SQL (or handling the resulting error). This includes verifying that complex
// identifiers are enclosed in double quotes, etc.
func GetTypeFromValidSQLSyntax(sql string) (tree.ResolvableTypeReference, error) {
	expr, err := ParseExpr(fmt.Sprintf("1::%s", sql))
	if err != nil {
		return nil, err
	}
	return GetTypeFromCastOrCollate(expr)
}

// GetTypeFromCastOrCollate returns the type of the given tree.Expr. The method
// assumes that the expression is either tree.CastExpr or tree.CollateExpr
// (which wraps the tree.CastExpr).
func GetTypeFromCastOrCollate(expr tree.Expr) (tree.ResolvableTypeReference, error) {
	// COLLATE clause has lower precedence than the cast, so if we have
	// something like `1::STRING COLLATE en`, it'll be parsed as
	// CollateExpr(CastExpr).
	if collate, ok := expr.(*tree.CollateExpr); ok {
		return types.MakeCollatedString(types.String, collate.Locale), nil
	}

	cast, ok := expr.(*tree.CastExpr)
	if !ok {
		return nil, errors.AssertionFailedf("expected a tree.CastExpr, but found %T", expr)
	}

	return cast.Type, nil
}

var errVarBitLengthNotPositive = pgerror.WithCandidateCode(
	errors.New("length for type varbit must be at least 1"), pgcode.InvalidParameterValue)

// newBitType creates a new BIT type with the given bit width.
func newBitType(width int32, varying bool) (*types.T, error) {
	if varying {
		if width < 1 {
			return nil, errVarBitLengthNotPositive
		}
		return types.MakeVarBit(width), nil
	}
	// The iconst32 pattern in the parser guarantees that the width is positive.
	return types.MakeBit(width), nil
}

var errFloatPrecAtLeast1 = pgerror.WithCandidateCode(
	errors.New("precision for type float must be at least 1 bit"), pgcode.InvalidParameterValue)
var errFloatPrecMax54 = pgerror.WithCandidateCode(
	errors.New("precision for type float must be less than 54 bits"), pgcode.InvalidParameterValue)

// newFloat creates a type for FLOAT with the given precision.
func newFloat(prec int64) (*types.T, error) {
	if prec < 1 {
		return nil, errFloatPrecAtLeast1
	}
	if prec <= 24 {
		return types.Float4, nil
	}
	if prec <= 54 {
		return types.Float, nil
	}
	return nil, errFloatPrecMax54
}

// newDecimal creates a type for DECIMAL with the given precision and scale.
func newDecimal(prec, scale int32) (*types.T, error) {
	if scale > prec {
		err := pgerror.WithCandidateCode(
			errors.Newf("scale (%d) must be between 0 and precision (%d)", scale, prec),
			pgcode.InvalidParameterValue)
		return nil, err
	}
	return types.MakeDecimal(prec, scale), nil
}

// arrayOf creates a type alias for an array of the given element type and fixed
// bounds. The bounds are currently ignored.
func arrayOf(
	ref tree.ResolvableTypeReference, bounds []int32,
) (tree.ResolvableTypeReference, error) {
	// If the reference is a statically known type, then return an array type,
	// rather than an array type reference.
	if typ, ok := tree.GetStaticallyKnownType(ref); ok {
		switch typ.Family() {
		case types.UnknownFamily, types.VoidFamily, types.TriggerFamily:
			// Do not allow arrays of these types. This is consistent with Postgres'
			// behavior.
			return nil, pgerror.Newf(pgcode.UndefinedObject, "type %s[] does not exist", typ.Name())
		}
		if err := types.CheckArrayElementType(typ); err != nil {
			return nil, err
		}
		return types.MakeArray(typ), nil
	}
	return &tree.ArrayTypeReference{ElementType: ref}, nil
}

// ParseDoBlockFn allows the SQL parser to parse a PL/pgSQL DO block body.
var ParseDoBlockFn func(tree.DoBlockOptions) (tree.DoBlockBody, error)
