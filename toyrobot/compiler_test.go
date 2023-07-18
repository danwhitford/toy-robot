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
					Type:   TOKEN_NUMBER,
					Value:  0,
					Lexeme: "0",
				},
				{
					Type:   TOKEN_NUMBER,
					Value:  0,
					Lexeme: "0",
				},
				{
					Type:   TOKEN_DIRECTION,
					Value:  NORTH,
					Lexeme: "NORTH",
				},
				{
					Type:   TOKEN_WORD,
					Value:  "PLACE",
					Lexeme: "PLACE",
				},
			},
			want: []byte{
				byte(OP_PUSH_VAL),
				byte(T_INT),
				byte(0),
				byte(OP_PUSH_VAL),
				byte(T_INT),
				byte(0),
				byte(OP_PUSH_VAL),
				byte(T_DIRECTION),
				byte(NORTH),
				byte(OP_EXEC_WORD),
				'P', 'L', 'A', 'C', 'E', 0,
			},
		},
		{
			input: []Token{
				{
					Type:   TOKEN_WORD,
					Value:  "MOVE",
					Lexeme: "MOVE",
				},
			},
			want: []byte{
				byte(OP_EXEC_WORD), 'M', 'O', 'V', 'E', 0,
			},
		},
		{
			input: []Token{
				{TOKEN_BOOL, true, "true"},
				{TOKEN_WORD, "IF", "IF"},
				{TOKEN_STRING, "hello", "\"hello\""},
				{TOKEN_WORD, ".", "."},
				{TOKEN_WORD, "FI", "FI"},
			},
			want: []byte{
				byte(OP_PUSH_VAL),
				byte(T_BOOL),
				byte(1),
				byte(OP_EXEC_WORD),
				'I', 'F', 0,
				19,
				byte(OP_PUSH_VAL),
				byte(T_STRING),
				'h', 'e', 'l', 'l', 'o', 0,
				byte(OP_EXEC_WORD),
				'.', 0,
			},
		},
		{
			input: []Token{
				{TOKEN_NUMBER, 5, "5"},
				{TOKEN_WORD, "DUP", "DUP"},
				{TOKEN_NUMBER, 5, "5"},
				{TOKEN_WORD, "EQ", "EQ"},
				{TOKEN_WORD, "IF", "IF"},
				{TOKEN_STRING, "5", "\"5\""},
				{TOKEN_WORD, ".", "."},
				{TOKEN_WORD, "ELSE", "ELSE"},
				{TOKEN_WORD, "DUP", "DUP"},
				{TOKEN_NUMBER, 5, "5"},
				{TOKEN_WORD, "GT", "GT"},
				{TOKEN_WORD, "IF", "IF"},
				{TOKEN_STRING, "BIGUN", "\"BIGUN\""},
				{TOKEN_WORD, ".", "."},
				{TOKEN_WORD, "ELSE", "ELSE"},
				{TOKEN_STRING, "SMALLUN", "\"SMALLUN\""},
				{TOKEN_WORD, ".", "."},
				{TOKEN_WORD, "FI", "FI"},
				{TOKEN_WORD, "FI", "FI"},
				{TOKEN_WORD, "DROP", "DROP"},
			},
			want: []byte{
				byte(OP_PUSH_VAL), // 0
				byte(T_INT),
				byte(5),
				byte(OP_EXEC_WORD),
				'D', 'U', 'P', 0,
				byte(OP_PUSH_VAL), // 8
				byte(T_INT),
				byte(5), // 10
				byte(OP_EXEC_WORD),
				'E', 'Q', 0, // 14
				byte(OP_EXEC_WORD), // 15
				'I', 'F', 0,
				0x21,
				byte(OP_PUSH_VAL),
				byte(T_STRING),
				'5', 0,
				byte(OP_EXEC_WORD),
				'.', 0,
				byte(OP_EXEC_WORD),
				'J', 'M', 'P', 0,
				0x50,
				byte(OP_EXEC_WORD),
				'D', 'U', 'P', 0,
				byte(OP_PUSH_VAL),
				byte(T_INT),
				byte(5), // 25
				byte(OP_EXEC_WORD),
				'G', 'T', 0,
				byte(OP_EXEC_WORD),
				'I', 'F', 0,
				0x43, // 30
				byte(OP_PUSH_VAL),
				byte(T_STRING),
				'B', 'I', 'G', 'U', 'N', 0,
				byte(OP_EXEC_WORD),
				'.', 0, // 35
				byte(OP_EXEC_WORD),
				'J', 'M', 'P', 0,
				0x50,
				byte(OP_PUSH_VAL),
				byte(T_STRING),
				'S', 'M', 'A', 'L', 'L', 'U', 'N', 0,
				byte(OP_EXEC_WORD),
				'.', 0,
				byte(OP_EXEC_WORD),
				'D', 'R', 'O', 'P', 0,
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
