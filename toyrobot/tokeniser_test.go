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
				{Type: TOKEN_WORD, Value: "PLACE", Lexeme: "PLACE"},
			},
		},
		{
			"RIGHT",
			[]Token{
				{Type: TOKEN_WORD, Value: "RIGHT", Lexeme: "RIGHT"},
			},
		},
		{
			"REPORT",
			[]Token{
				{Type: TOKEN_WORD, Value: "REPORT", Lexeme: "REPORT"},
			},
		},
		{
			"MOVE LEFT RIGHT REPORT",
			[]Token{
				{Type: TOKEN_WORD, Value: "MOVE", Lexeme: "MOVE"},
				{Type: TOKEN_WORD, Value: "LEFT", Lexeme: "LEFT"},
				{Type: TOKEN_WORD, Value: "RIGHT", Lexeme: "RIGHT"},
				{Type: TOKEN_WORD, Value: "REPORT", Lexeme: "REPORT"},
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
		{
			"10 20 30 40",
			[]Token{
				{Type: TOKEN_NUMBER, Value: 10, Lexeme: "10"},
				{Type: TOKEN_NUMBER, Value: 20, Lexeme: "20"},
				{Type: TOKEN_NUMBER, Value: 30, Lexeme: "30"},
				{Type: TOKEN_NUMBER, Value: 40, Lexeme: "40"},
			},
		},
		{
			"+ - * /",
			[]Token{
				{Type: TOKEN_WORD, Value: "+", Lexeme: "+"},
				{Type: TOKEN_WORD, Value: "-", Lexeme: "-"},
				{Type: TOKEN_WORD, Value: "*", Lexeme: "*"},
				{Type: TOKEN_WORD, Value: "/", Lexeme: "/"},
			},
		},
		{
			"\"hello world\"",
			[]Token{
				{Type: TOKEN_STRING, Value: "hello world", Lexeme: "\"hello world\""},
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
