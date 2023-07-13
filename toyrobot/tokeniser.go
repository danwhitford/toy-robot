package toyrobot

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/danwhitford/toyrobot/belt"
)

type RobotTokeniser struct {
	belt *belt.Belt[rune]
}

type TokenType byte

//go:generate stringer -type=TokenType
const (
	TOKEN_NUMBER TokenType = iota
	TOKEN_COMMA
	TOKEN_DIRECTION
	TOKEN_WORD
)

type Token struct {
	Type   TokenType
	Value  any
	Lexeme string
}

func (t *RobotTokeniser) Tokenise(input string) ([]Token, error) {
	tokens := make([]Token, 0)

	t.belt = belt.NewBelt[rune]([]rune(input))

	for t.belt.HasNext() {
		currentRune, err := t.belt.Peek()
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
		case currentRune == ',':
			token, err := t.getTokenComma()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case currentRune == '#':
			for currentRune != '\n' && t.belt.HasNext() {
				currentRune, err = t.belt.GetNext()
				if err != nil {
					return []Token{}, err
				}
			}
		case unicode.IsPrint(currentRune):
			token, err := t.getTokenAlpha()
			if err != nil {
				return []Token{}, err
			}
			tokens = append(tokens, token)
		case unicode.IsSpace(currentRune):
			t.belt.GetNext()
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

func (t *RobotTokeniser) getTokenComma() (Token, error) {
	curr, err := t.belt.GetNext()
	if err != nil {
		return Token{}, err
	}
	if curr == ',' {
		return Token{Type: TOKEN_COMMA, Value: nil, Lexeme: ","}, nil
	}

	return Token{}, fmt.Errorf("invalid token, expecting COMMA")
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
	default:
		return Token{Type: TOKEN_WORD, Value: strings.ToUpper(lexeme), Lexeme: lexeme}, nil
	}
}

func (t *RobotTokeniser) getLexeme() (string, error) {
	lexeme := ""
	for t.belt.HasNext() {
		curr, err := t.belt.GetNext()
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
