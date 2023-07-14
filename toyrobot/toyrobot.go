package toyrobot

import (
	"fmt"
	"io"
	"os"

	"github.com/danwhitford/toyrobot/belt"
	"github.com/danwhitford/toyrobot/stack"
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
			case T_BOOL:
				vi, err := r.Instructions.GetNext()
				if err != nil {
					return err
				}
				v := vi != 0
				r.RobotValueStack.Push(RobotValue{Type: t, Value: v})
			case T_STRING:
				vs := make([]byte, 0)
				vi, err := r.Instructions.GetNext()
				if err != nil {
					return err
				}
				for vi != 0 {
					vs = append(vs, vi)
					vi, err = r.Instructions.GetNext()
					if err != nil {
						return err
					}
				}
				v := string(vs)
				r.RobotValueStack.Push(RobotValue{Type: t, Value: v})
			default:
				return fmt.Errorf("invalid type %s", t)
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
