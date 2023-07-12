package toyrobot

import (
	"fmt"
	"github.com/danwhitford/toyrobot/belt"
)

//go:generate stringer -type=RobotType
type RobotType byte

const (
	T_INT RobotType = iota
	T_DIRECTION
)

//go:generate stringer -type=Instruction
type Instruction byte

const (
	OP_PLACE Instruction = iota
	OP_PUSH_VAL
	OP_POP_VAL
	OP_EXEC_WORD
)

type RobotCompiler struct {
	belt *belt.Belt[Token]
}

func (r *RobotCompiler) Compile(input []Token) ([]byte, error) {
	r.belt = belt.NewBelt[Token](input)

	instructions := make([]byte, 0)
	for r.belt.HasNext() {
		token, err := r.belt.GetNext()
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case TOKEN_PLACE:
			instructions = append(instructions, byte(OP_PLACE))
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
			bytes := append([]byte(tokenVal), 0)
			instructions = append(
				instructions,
				byte(OP_EXEC_WORD),
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
