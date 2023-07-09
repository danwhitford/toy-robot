package toyrobot

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenise(t *testing.T) {
	table := []struct {
		input    string
		expected []Token
	}{
		{
			"PLACE 3,2,NORTH",
			[]Token{
				{Type: TOKEN_PLACE, Value: nil, Lexeme: "PLACE"},
				{Type: TOKEN_NUMBER, Value: 3, Lexeme: "3"},
				{Type: TOKEN_COMMA, Value: nil, Lexeme: ","},
				{Type: TOKEN_NUMBER, Value: 2, Lexeme: "2"},
				{Type: TOKEN_COMMA, Value: nil, Lexeme: ","},
				{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: "NORTH"},
			},
		},
		{
			"RIGHT",
			[]Token{
				{Type: TOKEN_RIGHT, Value: nil, Lexeme: "RIGHT"},
			},
		},
		{
			"REPORT",
			[]Token{
				{Type: TOKEN_REPORT, Value: nil, Lexeme: "REPORT"},
			},
		},
		{
			"MOVE LEFT RIGHT REPORT",
			[]Token{
				{Type: TOKEN_MOVE, Value: nil, Lexeme: "MOVE"},
				{Type: TOKEN_LEFT, Value: nil, Lexeme: "LEFT"},
				{Type: TOKEN_RIGHT, Value: nil, Lexeme: "RIGHT"},
				{Type: TOKEN_REPORT, Value: nil, Lexeme: "REPORT"},
			},
		},
		{
			"NORTH SOUTH EAST WEST",
			[]Token{
				{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: "NORTH"},
				{Type: TOKEN_DIRECTION, Value: SOUTH, Lexeme: "SOUTH"},
				{Type: TOKEN_DIRECTION, Value: EAST, Lexeme: "EAST"},
				{Type: TOKEN_DIRECTION, Value: WEST, Lexeme: "WEST"},
			},
		},
	}

	for _, tst := range table {
		tokeniser := Tokeniser{}
		got, err := tokeniser.Tokenise(tst.input)
		if err != nil {
			t.Errorf("Error tokenising '%s': '%s'", tst.input, err)
		}

		if diff := cmp.Diff(tst.expected, got); diff != "" {
			t.Errorf("Tokenise(%s) mismatch (-want +got):\n%s", tst.input, diff)
		}
	}
}
