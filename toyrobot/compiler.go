package toyrobot

import (
	"fmt"

	"github.com/danwhitford/toyrobot/belt"
	"github.com/danwhitford/toyrobot/stack"
)

//go:generate stringer -type=Instruction
type Instruction byte

const (
	OP_PUSH_VAL Instruction = iota
	OP_EXEC_WORD
)

type IfFrame struct {
	Location int
}

type RobotCompiler struct {
	tokens  *belt.Belt[Token]
	ifStack stack.RobotStack[IfFrame]
}

func (r *RobotCompiler) Compile(input []Token) ([]byte, error) {
	r.tokens = belt.NewBelt[Token](input)

	instructions := make([]byte, 0)
	for r.tokens.HasNext() {
		token, err := r.tokens.GetNext()
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case TOKEN_NUMBER:
			instructions = append(
				instructions,
				byte(OP_PUSH_VAL),
				byte(T_INT),
				byte(token.Value.(int)),
			)
		case TOKEN_DIRECTION:
			instructions = append(
				instructions,
				byte(OP_PUSH_VAL),
				byte(T_DIRECTION),
				byte(token.Value.(Direction)),
			)
		case TOKEN_WORD:
			tokenVal, ok := token.Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid token value '%v'", token.Value)
			}
			switch tokenVal {
			case "IF":
				instructions = append(
					instructions,
					byte(OP_EXEC_WORD),
					'I', 'F', 0, 0, // placeholder for FI location
				)
				r.ifStack.Push(IfFrame{
					Location: len(instructions),
				})
			case "FI":
				ifFrame, err := r.ifStack.Pop()
				if err != nil {
					return nil, err
				}
				instructions[ifFrame.Location-1] = byte(len(instructions))
			default:
				bytes := append([]byte(tokenVal), 0)
				instructions = append(
					instructions,
					byte(OP_EXEC_WORD),
				)
				instructions = append(
					instructions,
					bytes...,
				)
			}
		case TOKEN_BOOL:
			boolVal, ok := token.Value.(bool)
			if !ok {
				return nil, fmt.Errorf("invalid token value '%v'", token.Value)
			}
			var byt byte
			if boolVal {
				byt = 1
			}

			instructions = append(
				instructions,
				byte(OP_PUSH_VAL),
				byte(T_BOOL),
				byt,
			)
		case TOKEN_STRING:
			tokenVal, ok := token.Value.(string)
			if !ok {
				return nil, fmt.Errorf("invalid token value '%v'", token.Value)
			}
			bytes := append([]byte(tokenVal), 0)
			instructions = append(
				instructions,
				byte(OP_PUSH_VAL),
				byte(T_STRING),
			)
			instructions = append(
				instructions,
				bytes...,
			)
		default:
			return nil, fmt.Errorf("invalid instruction '%v'", token)
		}
	}

	return instructions, nil
}
