package toyrobot

import (
	"fmt"
	"io"
	"os"

	"github.com/danwhitford/toyrobot/belt"
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

// This is a robot
type Robot struct {
	X, Y            int
	F               Direction
	Placed          bool
	Output          io.Writer
	RobotTokeniser  *RobotTokeniser
	RobotCompiler   *RobotCompiler
	RobotValueStack *stack.RobotStack[RobotValue]
	Dictionary      map[string]func() error
	Instructions    *belt.Belt[byte]
}

func NewRobot() *Robot {
	stack := make(stack.RobotStack[RobotValue], 0)

	dict := make(map[string]func() error)

	r := Robot{
		Output:          os.Stdout,
		RobotTokeniser:  &RobotTokeniser{},
		RobotCompiler:   &RobotCompiler{},
		RobotValueStack: &stack,
		Dictionary:      dict,
	}

	r.LoadEnv()
	return &r
}

func (r *Robot) LoadEnv() {
	r.Dictionary["BOARD"] = r.printBoard
	r.Dictionary["REPORT"] = r.report
	r.Dictionary["RIGHT"] = r.right
	r.Dictionary["LEFT"] = r.left
	r.Dictionary["MOVE"] = r.move
	r.Dictionary["PLACE"] = r.place
}

func (r *Robot) printBoard() error {
	hr := "+---+---+---+---+---+\n"
	cage := "| %s | %s | %s | %s | %s |\n"
	for y := 4; y >= 0; y-- {
		x := make([]interface{}, 5)
		for i := range x {
			x[i] = " "
		}
		if r.Placed && r.Y == y {
			switch r.F {
			case NORTH:
				x[r.X] = "^"
			case EAST:
				x[r.X] = ">"
			case SOUTH:
				x[r.X] = "v"
			case WEST:
				x[r.X] = "<"
			}
		}
		fmt.Fprint(r.Output, hr)
		fmt.Fprintf(r.Output, cage, x...)
	}
	fmt.Fprint(r.Output, hr)
	return nil
}

func (r *Robot) place() error {
	fv, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	yv, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	xv, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}

	f, ok := fv.Value.(Direction)
	if !ok {
		return fmt.Errorf("invalid direction %v", fv.Value)
	}
	y, ok := yv.Value.(int)
	if !ok {
		return fmt.Errorf("invalid y %v", yv.Value)
	}
	x, ok := xv.Value.(int)
	if !ok {
		return fmt.Errorf("invalid x %v", xv.Value)
	}

	if x < 0 || x > 4 || y < 0 || y > 4 {
		return nil
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

func (r *Robot) runInstructions() error {
	for r.Instructions.HasNext() {
		currentInstruction, err := r.Instructions.GetNext()
		if err != nil {
			return err
		}
		switch currentInstruction {
		case byte(OP_PUSH_VAL):
			typeInstruction, err := r.Instructions.GetNext()
			if err != nil {
				return err
			}
			t := RobotType(typeInstruction)
			switch t {
			case T_INT:
				vi, err := r.Instructions.GetNext()
				if err != nil {
					return err
				}
				v := int(vi)
				r.RobotValueStack.Push(RobotValue{Type: t, Value: v})
			case T_DIRECTION:
				vi, err := r.Instructions.GetNext()
				if err != nil {
					return err
				}
				v := Direction(vi)
				r.RobotValueStack.Push(RobotValue{Type: t, Value: v})
			}
		case byte(OP_EXEC_WORD):
			wordBytes := make([]byte, 0)
			t, err := r.Instructions.GetNext()
			if err != nil {
				return err
			}
			for t != 0 {
				wordBytes = append(wordBytes, t)
				t, err = r.Instructions.GetNext()
				if err != nil {
					return err
				}
			}
			word := string(wordBytes)
			fn, ok := r.Dictionary[word]
			if !ok {
				return fmt.Errorf("unknown word %s", word)
			}
			err = fn()
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid instruction %v", currentInstruction)
		}
	}
	return nil
}

func (r *Robot) RunProgram(instruction string) error {
	tokens, err := r.RobotTokeniser.Tokenise(instruction)
	if err != nil {
		return err
	}
	instructions, err := r.RobotCompiler.Compile(tokens)
	if err != nil {
		return err
	}
	r.Instructions = belt.NewBelt[byte](instructions)
	return r.runInstructions()
}
