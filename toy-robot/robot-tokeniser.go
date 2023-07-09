package toyrobot

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Tokeniser struct {
	input []rune
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

// TODO combine instructions into TOKEN_INSTRUCTION
func (t *Tokeniser) Tokenise(input string) ([]Token, error) {
	tokens := make([]Token, 0)

	t.input = []rune(input)
	t.ptr = 0
	t.size = len(input)

	for t.ptr < t.size {
		currentRune := t.getPeek()
		switch {
		case unicode.IsLetter(currentRune):
			token, err := t.getTokenAlpha()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case unicode.IsDigit(currentRune):
			token, err := t.getTokenNumber()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case currentRune == ',':
			token, err := t.getTokenComma()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

func (t *Tokeniser) getNext() rune {
	r := t.input[t.ptr]
	t.ptr++
	return r
}

func (t *Tokeniser) getPeek() rune {
	return t.input[t.ptr]
}

func (t *Tokeniser) getTokenNumber() (Token, error) {
	curr := t.getNext()
	number, err := strconv.Atoi(string(curr))
	if err != nil {
		return Token{}, fmt.Errorf("invalid token, expecting number but got '%s'", string(curr))
	}
	return Token{Type: TOKEN_NUMBER, Value: number, Lexeme: string(curr)}, nil
}

func (t *Tokeniser) getTokenComma() (Token, error) {
	if t.getNext() == ',' {
		return Token{Type: TOKEN_COMMA, Value: nil, Lexeme: ","}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting COMMA")
}

func (t *Tokeniser) getTokenAlpha() (Token, error) {
	lexeme := t.getLexeme()
	switch strings.ToUpper(lexeme) {
	case "PLACE":
		return Token{Type: TOKEN_PLACE, Value: nil, Lexeme: lexeme}, nil
	case "MOVE":
		return Token{Type: TOKEN_MOVE, Value: nil, Lexeme: lexeme}, nil
	case "LEFT":
		return Token{Type: TOKEN_LEFT, Value: nil, Lexeme: lexeme}, nil
	case "RIGHT":
		return Token{Type: TOKEN_RIGHT, Value: nil, Lexeme: lexeme}, nil
	case "REPORT":
		return Token{Type: TOKEN_REPORT, Value: nil, Lexeme: lexeme}, nil
	case "NORTH":
		return Token{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: lexeme}, nil
	case "EAST":
		return Token{Type: TOKEN_DIRECTION, Value: EAST, Lexeme: lexeme}, nil
	case "SOUTH":
		return Token{Type: TOKEN_DIRECTION, Value: SOUTH, Lexeme: lexeme}, nil
	case "WEST":
		return Token{Type: TOKEN_DIRECTION, Value: WEST, Lexeme: lexeme}, nil
	default:
		return Token{}, fmt.Errorf("invalid token, expecting instruction but got %s", lexeme)
	}
}

func (t *Tokeniser) getLexeme() string {
	lexeme := ""
	for t.ptr < t.size {
		currentRune := t.getNext()
		if unicode.IsLetter(currentRune) {
			lexeme += string(currentRune)
		} else {
			break
		}
	}
	return lexeme
}
