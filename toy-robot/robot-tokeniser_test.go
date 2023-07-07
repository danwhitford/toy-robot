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
				Token{Type: TOKEN_PLACE, Value: nil, Lexeme: "PLACE"},
				Token{Type: TOKEN_NUMBER, Value: 3, Lexeme: "3"},
				Token{Type: TOKEN_COMMA, Value: nil, Lexeme: ","},
				Token{Type: TOKEN_NUMBER, Value: 2, Lexeme: "2"},
				Token{Type: TOKEN_COMMA, Value: nil, Lexeme: ","},
				Token{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: "NORTH"},
			},
		},
		{
			"RIGHT",
			[]Token{
				Token{Type: TOKEN_RIGHT, Value: nil, Lexeme: "RIGHT"},
			},
		},
		{
			"REPORT",
			[]Token{
				Token{Type: TOKEN_REPORT, Value: nil, Lexeme: "REPORT"},
			},
		},
		{
			"MOVE LEFT RIGHT REPORT",
			[]Token{
				Token{Type: TOKEN_MOVE, Value: nil, Lexeme: "MOVE"},
				Token{Type: TOKEN_LEFT, Value: nil, Lexeme: "LEFT"},
				Token{Type: TOKEN_RIGHT, Value: nil, Lexeme: "RIGHT"},
				Token{Type: TOKEN_REPORT, Value: nil, Lexeme: "REPORT"},
			},
		},
		{
			"NORTH SOUTH EAST WEST",
			[]Token{
				Token{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: "NORTH"},
				Token{Type: TOKEN_DIRECTION, Value: SOUTH, Lexeme: "SOUTH"},
				Token{Type: TOKEN_DIRECTION, Value: EAST, Lexeme: "EAST"},
				Token{Type: TOKEN_DIRECTION, Value: WEST, Lexeme: "WEST"},
			},
		},
	}

	for _, tst := range table {
		tokeniser := Tokeniser{}
		got, err := tokeniser.Tokenise(tst.input)
		if err != nil {
			t.Errorf("Error tokenising %s: %s", tst.input, err)
		}

		if diff := cmp.Diff(tst.expected, got); diff != "" {
			t.Errorf("Tokenise(%s) mismatch (-want +got):\n%s", tst.input, diff)
		}
	}
}
