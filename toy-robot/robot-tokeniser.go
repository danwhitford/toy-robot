package toyrobot

import (
	"fmt"
	"strconv"
)

type Tokeniser struct {
	input string
	ptr   int
	size  int
}

type TokenType byte

//go:generate stringer -type=TokenType
const (
	TOKEN_PLACE TokenType = iota
	TOKEN_MOVE
	TOKEN_LEFT
	TOKEN_RIGHT
	TOKEN_REPORT
	TOKEN_NUMBER
	TOKEN_COMMA
	TOKEN_DIRECTION
)

type Token struct {
	Type   TokenType
	Value  any
	Lexeme string
}

func (t *Tokeniser) Tokenise(input string) ([]Token, error) {
	tokens := make([]Token, 0)

	t.input = input
	t.ptr = 0
	t.size = len(input)

	for t.ptr < t.size {
		currentRune := t.getNextByte()
		switch currentRune {
		case 'P':
			token, err := t.getTokenPlace()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case 'M':
			token, err := t.getTokenMove()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case 'L':
			token, err := t.getTokenLeft()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case 'R':
			next := t.getNextByte()
			if next == 'I' {
				token, err := t.getTokenRight()
				if err != nil {
					return []Token{}, err
				}
				tokens = append(tokens, token)
			} else if next == 'E' {
				token, err := t.getTokenReport()
				if err != nil {
					return []Token{}, err
				}
				tokens = append(tokens, token)
			} else {
				return []Token{}, fmt.Errorf(
					"invalid token, expecting RIGHT or REPORT but got %s at position %d",
					string(currentRune),
					t.ptr,
				)
			}
		case 'N':
			token, err := t.getTokenNorth()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case 'S':
			token, err := t.getTokenSouth()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case 'E':
			token, err := t.getTokenEast()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case 'W':
			token, err := t.getTokenWest()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			token, err := t.getTokenNumber()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case ',':
			token, err := t.getTokenComma()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (t *Tokeniser) getNextByte() byte {
	r := t.input[t.ptr]
	t.ptr++
	return r
}

func (t *Tokeniser) getLastRead() byte {
	return t.input[t.ptr-1]
}

func (t *Tokeniser) getTokenPlace() (Token, error) {
	if t.getNextByte() == 'L' &&
		t.getNextByte() == 'A' &&
		t.getNextByte() == 'C' &&
		t.getNextByte() == 'E' {
		return Token{Type: TOKEN_PLACE, Value: nil, Lexeme: "PLACE"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting PLACE")
}

func (t *Tokeniser) getTokenMove() (Token, error) {
	if t.getNextByte() == 'O' &&
		t.getNextByte() == 'V' &&
		t.getNextByte() == 'E' {
		return Token{Type: TOKEN_MOVE, Value: nil, Lexeme: "MOVE"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting MOVE")
}

func (t *Tokeniser) getTokenLeft() (Token, error) {
	if t.getNextByte() == 'E' &&
		t.getNextByte() == 'F' &&
		t.getNextByte() == 'T' {
		return Token{Type: TOKEN_LEFT, Value: nil, Lexeme: "LEFT"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting LEFT")
}

func (t *Tokeniser) getTokenRight() (Token, error) {
	if t.getNextByte() == 'G' &&
		t.getNextByte() == 'H' &&
		t.getNextByte() == 'T' {
		return Token{Type: TOKEN_RIGHT, Value: nil, Lexeme: "RIGHT"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting RIGHT")
}

func (t *Tokeniser) getTokenReport() (Token, error) {
	if t.getNextByte() == 'P' &&
		t.getNextByte() == 'O' &&
		t.getNextByte() == 'R' &&
		t.getNextByte() == 'T' {
		return Token{Type: TOKEN_REPORT, Value: nil, Lexeme: "REPORT"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting REPORT")
}

func (t *Tokeniser) getTokenNorth() (Token, error) {
	if t.getNextByte() == 'O' &&
		t.getNextByte() == 'R' &&
		t.getNextByte() == 'T' &&
		t.getNextByte() == 'H' {
		return Token{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: "NORTH"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting NORTH")
}

func (t *Tokeniser) getTokenSouth() (Token, error) {
	if t.getNextByte() == 'O' &&
		t.getNextByte() == 'U' &&
		t.getNextByte() == 'T' &&
		t.getNextByte() == 'H' {
		return Token{Type: TOKEN_DIRECTION, Value: SOUTH, Lexeme: "SOUTH"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting SOUTH")
}

func (t *Tokeniser) getTokenEast() (Token, error) {
	if t.getNextByte() == 'A' &&
		t.getNextByte() == 'S' &&
		t.getNextByte() == 'T' {
		return Token{Type: TOKEN_DIRECTION, Value: EAST, Lexeme: "EAST"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting EAST")
}

func (t *Tokeniser) getTokenWest() (Token, error) {
	if t.getNextByte() == 'E' &&
		t.getNextByte() == 'S' &&
		t.getNextByte() == 'T' {
		return Token{Type: TOKEN_DIRECTION, Value: WEST, Lexeme: "WEST"}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting WEST")
}

func (t *Tokeniser) getTokenNumber() (Token, error) {
	curr := t.getLastRead()
	number, err := strconv.Atoi(string(curr))
	if err != nil {
		return Token{}, fmt.Errorf("invalid token, expecting number but got %s", string(curr))
	}
	return Token{Type: TOKEN_NUMBER, Value: number, Lexeme: string(curr)}, nil
}

func (t *Tokeniser) getTokenComma() (Token, error) {
	if t.getLastRead() == ',' {
		return Token{Type: TOKEN_COMMA, Value: nil, Lexeme: ","}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting COMMA")
}
