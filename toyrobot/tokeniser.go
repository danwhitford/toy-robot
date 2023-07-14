package toyrobot

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/danwhitford/toyrobot/belt"
)

type RobotTokeniser struct {
	input *belt.Belt[rune]
}

type TokenType byte

//go:generate stringer -type=TokenType
const (
	TOKEN_NUMBER TokenType = iota
	TOKEN_DIRECTION
	TOKEN_WORD
	TOKEN_BOOL
	TOKEN_STRING
)

type Token struct {
	Type   TokenType
	Value  any
	Lexeme string
}

func (t *RobotTokeniser) Tokenise(input string) ([]Token, error) {
	tokens := make([]Token, 0)

	t.input = belt.NewBelt[rune]([]rune(input))

	for t.input.HasNext() {
		currentRune, err := t.input.Peek()
		if err != nil {
			return []Token{}, err
		}
		switch {
		case unicode.IsDigit(currentRune):
			token, err := t.getTokenNumber()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case currentRune == '#':
			for currentRune != '\n' && t.input.HasNext() {
				currentRune, err = t.input.GetNext()
				if err != nil {
					return []Token{}, err
				}
			}
		case currentRune == '"':
			token, err := t.getTokenString()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case unicode.IsPrint(currentRune):
			token, err := t.getTokenAlpha()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case unicode.IsSpace(currentRune):
			t.input.GetNext()
		default:
			return []Token{}, fmt.Errorf("invalid token, unexpected '%s'", string(currentRune))
		}
	}

	return tokens, nil
}

func (t *RobotTokeniser) getTokenNumber() (Token, error) {
	lexeme, err := t.getLexeme()
	if err != nil {
		return Token{}, err
	}
	number, err := strconv.Atoi(lexeme)
	if err != nil {
		return Token{}, fmt.Errorf("invalid token, expecting number but got '%s'", string(lexeme))
	}
	return Token{Type: TOKEN_NUMBER, Value: number, Lexeme: lexeme}, nil
}

func (t *RobotTokeniser) getTokenAlpha() (Token, error) {
	lexeme, err := t.getLexeme()
	if err != nil {
		return Token{}, err
	}
	switch strings.ToUpper(lexeme) {
	case "NORTH":
		return Token{Type: TOKEN_DIRECTION, Value: NORTH, Lexeme: lexeme}, nil
	case "EAST":
		return Token{Type: TOKEN_DIRECTION, Value: EAST, Lexeme: lexeme}, nil
	case "SOUTH":
		return Token{Type: TOKEN_DIRECTION, Value: SOUTH, Lexeme: lexeme}, nil
	case "WEST":
		return Token{Type: TOKEN_DIRECTION, Value: WEST, Lexeme: lexeme}, nil
	case "TRUE":
		return Token{Type: TOKEN_BOOL, Value: true, Lexeme: lexeme}, nil
	case "FALSE":
		return Token{Type: TOKEN_BOOL, Value: false, Lexeme: lexeme}, nil
	default:
		return Token{Type: TOKEN_WORD, Value: strings.ToUpper(lexeme), Lexeme: lexeme}, nil
	}
}

func (t *RobotTokeniser) getLexeme() (string, error) {
	lexeme := ""
	for t.input.HasNext() {
		curr, err := t.input.GetNext()
		if err != nil {
			return "", err
		}
		currentRune := curr
		if !unicode.IsSpace(currentRune) {
			lexeme += string(currentRune)
		} else {
			break
		}
	}
	return lexeme, nil
}

func (t *RobotTokeniser) getTokenString() (Token, error) {
	_, err := t.input.GetNext()
	if err != nil {
		return Token{}, err
	}
	lexeme := ""
	for t.input.HasNext() {
		curr, err := t.input.GetNext()
		if err != nil {
			return Token{}, err
		}
		currentRune := curr
		if currentRune == '"' {
			break
		} else {
			lexeme += string(currentRune)
		}
	}
	return Token{TOKEN_STRING, lexeme, fmt.Sprintf("\"%s\"", lexeme)}, nil
}
