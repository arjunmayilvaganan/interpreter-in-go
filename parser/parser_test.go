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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil!")
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Number of statements expected=%d, got=%d",
				1, len(program.Statements))
		}
		statement := program.Statements[0]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
			return
		}

		val := statement.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Number of statements expected=%d, got=%d",
				1, len(program.Statements))
		}

		statement := program.Statements[0]
		returnStatement, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("statement is expected=%s, got=%T", "*ast.ReturnStatement", statement)
		}
		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("Statement TokenLiteral expected=%s, got=%s",
				"return", returnStatement.TokenLiteral())
		}
		if testLiteralExpression(t, returnStatement.ReturnValue, tt.expectedValue) {
			return
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
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
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
		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
	}

	for _, tt := range infixTests {
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

		if !testInfixExpression(t, s.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			t.Fail()
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b",
			"((-a) * b)",
		},
		{"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%s, got=%s", tt.expected, actual)
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

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Unexpected expression type. expected=%s, got=%t", "*ast.Identifier", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("Unexpected identifier value. expected=%s, got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Unexpected token literal. expected=%s, got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	}
	t.Errorf("Received unexpected type. Got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	expression, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("s is expected=%s, got=%T", "*ast.InfixExpression", expression)
	}
	if expression.Operator != operator {
		t.Errorf("expression.Operator expected=%s, got=%s", operator, expression.Operator)
	}
	if !testLiteralExpression(t, expression.Left, left) {
		return false
	}
	if !testLiteralExpression(t, expression.Right, right) {
		return false
	}
	return true
}
