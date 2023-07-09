package toyrobot

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompile(t *testing.T) {
	table := []struct {
		input []Token
		want  []byte
	}{
		{
			input: []Token{
				{
					Type:   TOKEN_PLACE,
					Value:  nil,
					Lexeme: "PLACE",
				},
				{
					Type:   TOKEN_NUMBER,
					Value:  0,
					Lexeme: "0",
				},
				{
					Type:   TOKEN_COMMA,
					Value:  nil,
					Lexeme: ",",
				},
				{
					Type:   TOKEN_NUMBER,
					Value:  0,
					Lexeme: "0",
				},
				{
					Type:   TOKEN_COMMA,
					Value:  nil,
					Lexeme: ",",
				},
				{
					Type:   TOKEN_DIRECTION,
					Value:  NORTH,
					Lexeme: "NORTH",
				},
			},
			want: []byte{
				byte(PLACE), byte(0), byte(0), byte(NORTH),
			},
		},
		{
			input: []Token{
				{
					Type:   TOKEN_MOVE,
					Value:  nil,
					Lexeme: "MOVE",
				},
			},
			want: []byte{
				0x01,
			},
		},
	}

	for _, test := range table {
		compiler := RobotCompiler{}
		got, err := compiler.Compile(test.input)
		if err != nil {
			t.Errorf("Compile(%v) returned error %v", test.input, err)
		}
		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("Compile(%v) returned unexpected diff (-want +got):\n%s", test.input, diff)
		}
	}
}