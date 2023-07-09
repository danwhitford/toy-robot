package toyrobot

import (
	"fmt"
	"io"
	"os"
)

//go:generate stringer -type=Direction
type Direction byte

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)

type Instruction byte

const (
	PLACE Instruction = iota
	MOVE
	LEFT
	RIGHT
	REPORT
)

// This is a robot
type Robot struct {
	X, Y           int
	F              Direction
	Placed         bool
	Output         io.Writer
	RobotTokeniser *Tokeniser
	RobotCompiler  *RobotCompiler
}

func NewRobot() *Robot {
	return &Robot{
		Output:         os.Stdout,
		RobotTokeniser: &Tokeniser{},
		RobotCompiler:  &RobotCompiler{},
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

// TODO make this better
func (r *Robot) runInstructions(instructions []byte) error {
	idx := 0
	for idx < len(instructions) {
		fmt.Printf("idx %d, instruction %v of %d\n", idx, instructions[idx], len(instructions))
		switch instructions[idx] {
		case byte(PLACE):
			idx++
			x := int(instructions[idx])
			idx++
			y := int(instructions[idx])
			idx++
			f := Direction(instructions[idx])
			idx++
			err := r.place(x, y, f)
			if err != nil {
				return err
			}
		case byte(MOVE):
			idx++
			err := r.move()
			if err != nil {
				return err
			}
		case byte(LEFT):
			idx++
			err := r.left()
			if err != nil {
				return err
			}
		case byte(RIGHT):
			idx++
			err := r.right()
			if err != nil {
				return err
			}
		case byte(REPORT):
			idx++
			err := r.report()
			if err != nil {
				return err
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
