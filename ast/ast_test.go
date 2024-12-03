package ast

import (
	"nibbl/token"
	"testing"
)

func TestString(t *testing.T) {
	expected := "let myVar = anotherVar;"
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != expected {
		t.Errorf("program.String, expected=\n%s\ngot=\n%s", expected, program.String())
	}
}
