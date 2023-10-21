package main

import (
	"go/ast"
)

// Given an AST node corresponding to a function (guaranteed to be
// of the form `func f(x, y byte) byte`), compile it into assembly
// code.
//
// Recall from the README that the input parameters `x` and `y` should
// be read from memory addresses `1` and `2`, and the return value
// should be written to memory address `0`.
func compile(fn *ast.FuncDecl) (string, error) {
    for _, stmt := range fn.Body.List {
        if ret, ok := stmt.(*ast.ReturnStmt); ok {
            expr := ret.Results[0]
            if len(ret.Results) == 1 {
                if basicLit, ok := expr.(*ast.BasicLit); ok {
                    return "pushi " + basicLit.Value + `
                    pop 0
                    halt`, nil
                }
            }
        }
    }
	// TODO
	return "halt\n", nil
}
