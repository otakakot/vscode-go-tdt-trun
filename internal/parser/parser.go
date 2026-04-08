package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

type SubTest struct {
	Func string `json:"func"`
	Name string `json:"name"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func ExtractSubTests(filename string) ([]SubTest, error) {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, absPath, nil, 0)
	if err != nil {
		return nil, err
	}

	var subtests []SubTest

	for _, decl := range f.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || !isTestFunc(funcDecl) {
			continue
		}

		subs := extractFromFunc(fset, absPath, funcDecl)
		subtests = append(subtests, subs...)
	}

	return subtests, nil
}

func isTestFunc(fn *ast.FuncDecl) bool {
	if !strings.HasPrefix(fn.Name.Name, "Test") {
		return false
	}

	params := fn.Type.Params
	if params == nil || len(params.List) != 1 {
		return false
	}

	return isTestingTType(params.List[0].Type)
}

func isTestingTType(expr ast.Expr) bool {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}

	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "testing" && sel.Sel.Name == "T"
}

func tParamName(fn *ast.FuncDecl) string {
	names := fn.Type.Params.List[0].Names
	if len(names) > 0 {
		return names[0].Name
	}

	return ""
}

type rangeInfo struct {
	keyVar      string
	valVar      string
	rangeTarget string
	rangeExpr   ast.Expr
	bodyStart   token.Pos
	bodyEnd     token.Pos
}

func extractFromFunc(fset *token.FileSet, filename string, funcDecl *ast.FuncDecl) []SubTest {
	testName := funcDecl.Name.Name

	tParam := tParamName(funcDecl)
	if tParam == "" {
		return nil
	}

	body := funcDecl.Body
	assignments := collectAssignments(body)

	var ranges []rangeInfo

	ast.Inspect(body, func(n ast.Node) bool {
		rs, ok := n.(*ast.RangeStmt)
		if !ok {
			return true
		}

		ri := rangeInfo{
			bodyStart: rs.Body.Pos(),
			bodyEnd:   rs.Body.End(),
			rangeExpr: rs.X,
		}
		if rs.Key != nil {
			if ident, ok := rs.Key.(*ast.Ident); ok {
				ri.keyVar = ident.Name
			}
		}

		if rs.Value != nil {
			if ident, ok := rs.Value.(*ast.Ident); ok {
				ri.valVar = ident.Name
			}
		}

		if ident, ok := rs.X.(*ast.Ident); ok {
			ri.rangeTarget = ident.Name
		}

		ranges = append(ranges, ri)

		return true
	})

	var subTests []SubTest

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok || len(call.Args) < 2 {
			return true
		}

		if !isTRunCall(call, tParam) {
			return true
		}

		var info *rangeInfo

		for i := range ranges {
			ri := &ranges[i]
			if call.Pos() > ri.bodyStart && call.End() < ri.bodyEnd {
				if info == nil || (ri.bodyEnd-ri.bodyStart < info.bodyEnd-info.bodyStart) {
					info = ri
				}
			}
		}

		subs := resolveNames(fset, filename, testName, call.Args[0], info, assignments)

		subTests = append(subTests, subs...)

		return true
	})

	return subTests
}

func collectAssignments(body *ast.BlockStmt) map[string]ast.Expr {
	assignments := make(map[string]ast.Expr)

	ast.Inspect(body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			for i, lhs := range stmt.Lhs {
				if ident, ok := lhs.(*ast.Ident); ok && i < len(stmt.Rhs) {
					assignments[ident.Name] = stmt.Rhs[i]
				}
			}
		case *ast.DeclStmt:
			genDecl, ok := stmt.Decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				break
			}

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				for i, name := range valueSpec.Names {
					if i < len(valueSpec.Values) {
						assignments[name.Name] = valueSpec.Values[i]
					}
				}
			}
		}

		return true
	})

	return assignments
}

func isTRunCall(call *ast.CallExpr, tParam string) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Run" {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == tParam
}

func resolveNames(fset *token.FileSet, filename, testFuncName string, arg ast.Expr, ctx *rangeInfo, assignments map[string]ast.Expr) []SubTest {
	switch a := arg.(type) {
	// case *ast.BasicLit:
	// 	// Pattern: t.Run("literal name", ...) — not table-driven, intentionally skipped.
	// 	// The official VS Code Go extension already supports running/debugging literal subtests
	// 	// via the "Go: Subtest At Cursor" command.
	// 	// See: https://github.com/golang/vscode-go/wiki/commands
	case *ast.SelectorExpr:
		// Pattern: t.Run(tt.name, ...) — slice table-driven
		if ctx == nil {
			return nil
		}

		ident, ok := a.X.(*ast.Ident)
		if !ok || ident.Name != ctx.valVar {
			return nil
		}

		fieldName := a.Sel.Name

		expr := resolveRangeExpr(ctx, assignments)
		if expr == nil {
			return nil
		}

		return extractFieldValues(fset, filename, testFuncName, expr, fieldName)
	case *ast.Ident:
		// Pattern: t.Run(name, ...) — map table-driven
		if ctx == nil {
			return nil
		}

		if a.Name != ctx.keyVar {
			return nil
		}

		expr := resolveRangeExpr(ctx, assignments)
		if expr == nil {
			return nil
		}

		return extractMapKeys(fset, filename, testFuncName, expr)
	}

	return nil
}

func resolveRangeExpr(ctx *rangeInfo, assignments map[string]ast.Expr) ast.Expr {
	if ctx.rangeTarget != "" {
		if expr, ok := assignments[ctx.rangeTarget]; ok {
			return expr
		}
	}

	if _, ok := ctx.rangeExpr.(*ast.CompositeLit); ok {
		return ctx.rangeExpr
	}

	return nil
}

func extractFieldValues(fset *token.FileSet, filename, testFuncName string, expr ast.Expr, fieldName string) []SubTest {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	fieldIndex := resolveFieldIndex(compLit, fieldName)

	var subtests []SubTest

	for _, elt := range compLit.Elts {
		innerLit, ok := elt.(*ast.CompositeLit)
		if !ok {
			continue
		}

		name, found := extractFieldFromLiteral(innerLit, fieldName, fieldIndex)
		if !found {
			continue
		}

		subtests = append(subtests, SubTest{
			Func: testFuncName,
			Name: name,
			File: filename,
			Line: fset.Position(innerLit.Pos()).Line,
		})
	}

	return subtests
}

// resolveFieldIndex determines the positional index of a named field in an anonymous struct slice type.
func resolveFieldIndex(compLit *ast.CompositeLit, fieldName string) int {
	arrayType, ok := compLit.Type.(*ast.ArrayType)
	if !ok {
		return -1
	}

	structType, ok := arrayType.Elt.(*ast.StructType)
	if !ok {
		return -1
	}

	idx := 0

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			if name.Name == fieldName {
				return idx
			}

			idx++
		}
	}

	return -1
}

func extractFieldFromLiteral(lit *ast.CompositeLit, fieldName string, fieldIndex int) (string, bool) {
	if len(lit.Elts) == 0 {
		return "", false
	}

	// Key-value syntax: {name: "test1", a: 1}
	if _, isKV := lit.Elts[0].(*ast.KeyValueExpr); isKV {
		for _, elt := range lit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != fieldName {
				continue
			}

			val, ok := kv.Value.(*ast.BasicLit)
			if !ok || val.Kind != token.STRING {
				return "", false
			}

			name, err := strconv.Unquote(val.Value)
			if err != nil {
				return "", false
			}

			return name, true
		}

		return "", false
	}

	// Positional syntax: {"test1", 1}
	if fieldIndex >= 0 && fieldIndex < len(lit.Elts) {
		val, ok := lit.Elts[fieldIndex].(*ast.BasicLit)
		if ok && val.Kind == token.STRING {
			name, err := strconv.Unquote(val.Value)
			if err != nil {
				return "", false
			}

			return name, true
		}
	}

	return "", false
}

func extractMapKeys(fset *token.FileSet, filename, testFuncName string, expr ast.Expr) []SubTest {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	var subTests []SubTest

	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			continue
		}

		name, err := strconv.Unquote(keyLit.Value)
		if err != nil {
			continue
		}

		subTests = append(subTests, SubTest{
			Func: testFuncName,
			Name: name,
			File: filename,
			Line: fset.Position(kv.Pos()).Line,
		})
	}

	return subTests
}
