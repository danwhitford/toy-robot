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
			"3 2 NORTH PLACE",
			[]Token{
				{Type: TOKEN_NUMBER, Value: 3, Lexeme: "3"},
				{Type: TOKEN_NUMBER, Value: 2, Lexeme: "2"},
				{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: "NORTH"},
				{Type: TOKEN_PLACE, Value: nil, Lexeme: "PLACE"},
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
		tokeniser := RobotTokeniser{}
		got, err := tokeniser.Tokenise(tst.input)
		if err != nil {
			t.Errorf("Error tokenising '%s': '%s'", tst.input, err)
		}

		if diff := cmp.Diff(tst.expected, got); diff != "" {
			t.Errorf("Tokenise(%s) mismatch (-want +got):\n%s", tst.input, diff)
		}
	}
}
