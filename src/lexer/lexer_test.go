package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {

	tests := []string{
		//"courses?filter=equals(",
		//"students?filter=equals(displayName,",
		//"courses?filter=equals(displayName,)",
		//"students?filter=equals(displayName,null)",
		"courses?filter=equals(displayName,'Brian Connor')",
		//"teachers?filter=equals(displayName,lastName)",
	}

	for _, test := range tests {
		l := New(test)

		for {
			tok := l.NextToken()
			tok.Print()
			if tok.Type == Eof {
				break
			}
		}
	}
}
