package toyrobot

import (
	"fmt"
	"github.com/danwhitford/toyrobot/belt"
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
		case TOKEN_MOVE:
			instructions = append(instructions, byte(OP_MOVE))
		case TOKEN_LEFT:
			instructions = append(instructions, byte(OP_LEFT))
		case TOKEN_RIGHT:
			instructions = append(instructions, byte(OP_RIGHT))
		case TOKEN_REPORT:
			instructions = append(instructions, byte(OP_REPORT))
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
		default:
			return nil, fmt.Errorf("invalid instruction '%v'", token)
		}
	}

	return instructions, nil
}
