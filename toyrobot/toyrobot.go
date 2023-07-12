package toyrobot

import (
	"fmt"
	"io"
	"os"

	"github.com/danwhitford/toyrobot/stack"
)

//go:generate stringer -type=Direction
type Direction byte

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

//go:generate stringer -type=Instruction
type Instruction byte

const (
	OP_PLACE Instruction = iota
	OP_MOVE
	OP_LEFT
	OP_RIGHT
	OP_REPORT
	OP_PUSH_VAL
	OP_POP_VAL
)

//go:generate stringer -type=RobotType
type RobotType byte

const (
	T_INT RobotType = iota
	T_DIRECTION
)

// This is a robot
type Robot struct {
	X, Y            int
	F               Direction
	Placed          bool
	Output          io.Writer
	RobotTokeniser  *RobotTokeniser
	RobotCompiler   *RobotCompiler
	RobotValueStack *stack.RobotStack[RobotValue]
}

func NewRobot() *Robot {
	stack := make(stack.RobotStack[RobotValue], 0)
	return &Robot{
		Output:          os.Stdout,
		RobotTokeniser:  &RobotTokeniser{},
		RobotCompiler:   &RobotCompiler{},
		RobotValueStack: &stack,
	}
}

func (r *Robot) place(x, y int, f Direction) error {
	if x < 0 || x > 4 || y < 0 || y > 4 {
		return fmt.Errorf("invalid coordinates %d,%d", x, y)
	}
	if f < NORTH || f > WEST {
		return fmt.Errorf("invalid facing %v", f)
	}

	r.X = x
	r.Y = y
	r.F = f
	r.Placed = true
	return nil
}

func (r *Robot) move() error {
	if !r.Placed {
		return nil
	}

	switch r.F {
	case NORTH:
		if r.Y < 4 {
			r.Y++
		}
	case EAST:
		if r.X < 4 {
			r.X++
		}
	case SOUTH:
		if r.Y > 0 {
			r.Y--
		}
	case WEST:
		if r.X > 0 {
			r.X--
		}
	}
	return nil
}

// Implement LEFT
func (r *Robot) left() error {
	if !r.Placed {
		return nil
	}
	switch r.F {
	case NORTH:
		r.F = WEST
	case EAST:
		r.F = NORTH
	case SOUTH:
		r.F = EAST
	case WEST:
		r.F = SOUTH
	}
	return nil
}

// Implement RIGHT
func (r *Robot) right() error {
	if !r.Placed {
		return nil
	}
	switch r.F {
	case NORTH:
		r.F = EAST
	case EAST:
		r.F = SOUTH
	case SOUTH:
		r.F = WEST
	case WEST:
		r.F = NORTH
	}
	return nil
}

// Implement REPORT
func (r *Robot) report() error {
	if !r.Placed {
		fmt.Fprintln(r.Output, "Robot not placed")
		return nil
	}
	fmt.Fprintf(r.Output, "%d,%d,%s\n", r.X, r.Y, r.F)
	return nil
}

// TODO make this better with instruction belt
func (r *Robot) runInstructions(instructions []byte) error {
	idx := 0
	for idx < len(instructions) {
		switch instructions[idx] {
		case byte(OP_PLACE):
			idx++
			dir, err := r.RobotValueStack.Pop()
			if err != nil {
				return err
			}
			y, err := r.RobotValueStack.Pop()
			if err != nil {
				return err
			}
			x, err := r.RobotValueStack.Pop()
			if err != nil {
				return err
			}
			r.place(x.Value.(int), y.Value.(int), dir.Value.(Direction))
		case byte(OP_MOVE):
			idx++
			err := r.move()
			if err != nil {
				return err
			}
		case byte(OP_LEFT):
			idx++
			err := r.left()
			if err != nil {
				return err
			}
		case byte(OP_RIGHT):
			idx++
			err := r.right()
			if err != nil {
				return err
			}
		case byte(OP_REPORT):
			idx++
			err := r.report()
			if err != nil {
				return err
			}
		case byte(OP_PUSH_VAL):
			idx++
			t := RobotType(instructions[idx])
			idx++
			switch t {
			case T_INT:
				v := int(instructions[idx])
				idx++
				r.RobotValueStack.Push(RobotValue{Type: t, Value: v})
			case T_DIRECTION:
				v := Direction(instructions[idx])
				idx++
				r.RobotValueStack.Push(RobotValue{Type: t, Value: v})
			}
		default:
			return fmt.Errorf("invalid instruction %v", instructions[idx])
		}
	}
	return nil
}

func (r *Robot) ReadInstruction(instruction string) error {
	tokens, err := r.RobotTokeniser.Tokenise(instruction)
	if err != nil {
		return err
	}
	instructions, err := r.RobotCompiler.Compile(tokens)
	if err != nil {
		return err
	}
	return r.runInstructions(instructions)
}
