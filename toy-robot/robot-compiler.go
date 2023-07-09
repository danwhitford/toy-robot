package toyrobot

import "fmt"

type RobotCompiler struct {
	belt *Belt[Token]
}

func (r *RobotCompiler) Compile(input []Token) ([]byte, error) {
	r.belt = NewBelt[Token](input)

	instructions := make([]byte, 0)
	for r.belt.HasNext() {
		token, err := r.belt.GetNext()
		if err != nil {
			return nil, err
		}
		switch token.Type {
		case TOKEN_PLACE:
			x, err := r.belt.GetNext()
			if err != nil {
				return nil, fmt.Errorf("invalid PLACE instruction: %v", err)
			}
			if x.Type != TOKEN_NUMBER {
				return nil, fmt.Errorf("invalid x coordinate %+v", x)
			}
			c, err := r.belt.GetNext()
			if err != nil {
				return nil, fmt.Errorf("invalid PLACE instruction: %v", err)
			}
			if c.Type != TOKEN_COMMA {
				return nil, fmt.Errorf("invalid PLACE instruction: %v", c)
			}
			y, err := r.belt.GetNext()
			if err != nil {
				return nil, fmt.Errorf("invalid PLACE instruction: %v", err)
			}
			if y.Type != TOKEN_NUMBER {
				return nil, fmt.Errorf("invalid y coordinate '%+v'", y)
			}
			c, err = r.belt.GetNext()
			if err != nil {
				return nil, fmt.Errorf("invalid PLACE instruction")
			}
			if c.Type != TOKEN_COMMA {
				return nil, fmt.Errorf("invalid PLACE instruction")
			}
			f, err := r.belt.GetNext()
			if err != nil {
				return nil, fmt.Errorf("invalid PLACE instruction")
			}
			if f.Type != TOKEN_DIRECTION {
				return nil, fmt.Errorf("invalid direction %+v", f)
			}
			instructions = append(instructions, byte(PLACE))
			instructions = append(instructions, byte(x.Value.(int)))
			instructions = append(instructions, byte(y.Value.(int)))
			instructions = append(instructions, byte(f.Value.(Direction)))
		case TOKEN_MOVE:
			instructions = append(instructions, byte(MOVE))
		case TOKEN_LEFT:
			instructions = append(instructions, byte(LEFT))
		case TOKEN_RIGHT:
			instructions = append(instructions, byte(RIGHT))
		case TOKEN_REPORT:
			instructions = append(instructions, byte(REPORT))
		case TOKEN_NUMBER:
			instructions = append(instructions, byte(token.Value.(int)))
		default:
			return nil, fmt.Errorf("invalid instruction '%v'", token)
		}
	}

	return instructions, nil
}
