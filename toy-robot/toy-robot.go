package toyrobot

import (
	"fmt"
	"io"
	"os"
	"strings"
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

type Robot struct {
	X, Y   int
	F      Direction
	Placed bool
	Output io.Writer
}

func NewRobot() *Robot {
	return &Robot{
		Output: os.Stdout,
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

func (r *Robot) runInstructions(instructions []byte) error {
	idx := 0
	for idx < len(instructions) {
		switch instructions[idx] {
		case byte(PLACE):
			idx++
			x := int(instructions[idx])
			idx++
			y := int(instructions[idx])
			idx++
			f := Direction(instructions[idx])
			return r.place(x, y, f)
		case byte(MOVE):
			idx++
			return r.move()
		case byte(LEFT):
			idx++
			return r.left()
		case byte(RIGHT):
			idx++
			return r.right()
		case byte(REPORT):
			idx++
			return r.report()
		default:
			return fmt.Errorf("invalid instruction %v", instructions[idx])
		}
	}
	return nil
}

func compileLine(instruction string) ([]byte, error) {
	switch {
	case strings.HasPrefix(instruction, "PLACE"):
		var x, y int
		var face string
		_, err := fmt.Sscanf(instruction, "PLACE %d,%d,%s", &x, &y, &face)
		if err != nil {
			return nil, err
		}
		f, err := stringToFacing(face)
		if err != nil {
			return nil, err
		}
		return []byte{byte(PLACE), byte(x), byte(y), byte(f)}, nil
	case instruction == "MOVE":
		return []byte{byte(MOVE)}, nil
	case instruction == "LEFT":
		return []byte{byte(LEFT)}, nil
	case instruction == "RIGHT":
		return []byte{byte(RIGHT)}, nil
	case instruction == "REPORT":
		return []byte{byte(REPORT)}, nil
	default:
		return nil, fmt.Errorf("invalid instruction %s", instruction)
	}
}

func (r *Robot) ReadInstruction(instruction string) error {
	instructions, err := compileLine(instruction)
	if err != nil {
		return err
	}
	return r.runInstructions(instructions)
}

func stringToFacing(face string) (Direction, error) {
	switch face {
	case "NORTH":
		return NORTH, nil
	case "EAST":
		return EAST, nil
	case "SOUTH":
		return SOUTH, nil
	case "WEST":
		return WEST, nil
	}
	return 0, fmt.Errorf("invalid facing %s", face)
}
