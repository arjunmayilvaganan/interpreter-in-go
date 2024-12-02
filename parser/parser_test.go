package parser

import (
	"nibbl/ast"
	"nibbl/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil!")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Number of statements expected=%d, got=%d",
			3, len(program.Statements))
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("Statement TokenLiteral expected=%s, got=%s", "let", s.TokenLiteral())
		return false
	}

	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is expected=%s, got=%T", "*ast.LetStatement", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value expected=%s, got=%s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name.TokenLiteral() expected=%s, got=%s",
			name, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}
