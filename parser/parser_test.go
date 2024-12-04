package parser

import (
	"fmt"
	"github.com/arjunmayilvaganan/nibbl/ast"
	"github.com/arjunmayilvaganan/nibbl/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) != 0 {
		t.Errorf("Parse errors encountered: %d\n", len(errors))
		for _, msg := range errors {
			t.Error(msg)
		}
		t.FailNow()
	}
}

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
	checkParserErrors(t, p)
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

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return add(15);
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Number of statements expected=%d, got=%d",
			3, len(program.Statements))
	}

	for _, s := range program.Statements {
		returnStatement, ok := s.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("s is expected=%s, got=%T", "*ast.ReturnStatement", s)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("Statement TokenLiteral expected=%s, got=%s",
				"return", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Number of statements expected=%d, got=%d",
			1, len(program.Statements))
	}

	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("s is expected=%s, got=%T", "*ast.ExpressionStatement", s)
	}

	ident, ok := s.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("s is expected=%s, got=%T", "*ast.ExpressionStatement", s)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value expected=%s, got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() expected=%s, got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Number of statements expected=%d, got=%d",
			1, len(program.Statements))
	}

	s, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("s is expected=%s, got=%T", "*ast.ExpressionStatement", s)
	}

	literal, ok := s.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("literal is expected=%s, got=%T", "*ast.IntegerLiteral", s)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value expected=%d, got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() expected=%s, got=%s", "5", literal.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Number of statements expected=%d, got=%d",
				1, len(program.Statements))
		}

		s, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("s is expected=%s, got=%T", "*ast.ExpressionStatement", s)
		}

		expression, ok := s.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("s is expected=%s, got=%T", "*ast.PrefixExpression", s)
		}
		if expression.Operator != tt.operator {
			t.Errorf("expression.Operator expected=%s, got=%s", tt.operator, expression.Operator)
		}
		if !testIntegerLiteral(t, expression.Right, tt.integerValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il was expected=%s, got=%T", "*ast.IntegerLiteral", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value was expected=%d, got=%d", value, integ.Value)
		return false
	}
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral() was expected=%d, got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}
