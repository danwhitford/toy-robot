package toyrobot

import (
	"fmt"

	"github.com/danwhitford/toyrobot/stack"
)

//go:generate go run ../templates/generate_builtins.go
func (r *Robot) LoadEnv() {
	// Robot stuff
	r.Dictionary["BOARD"] = r.printBoard
	r.Dictionary["REPORT"] = r.report
	r.Dictionary["RIGHT"] = r.right
	r.Dictionary["LEFT"] = r.left
	r.Dictionary["MOVE"] = r.move
	r.Dictionary["PLACE"] = r.place

	// Stack stuff
	r.Dictionary["."] = r.prn
	r.Dictionary["DUP"] = r.dup
	r.Dictionary["V"] = r.v
	r.Dictionary["CR"] = r.cr
	r.Dictionary["DROP"] = r.drop
	r.Dictionary["SWAP"] = r.swap
	r.Dictionary["OVER"] = r.over
	r.Dictionary["ROT"] = r.rot
	r.Dictionary["XX"] = r.clear

	// Math stuff
	r.Dictionary["+"] = r.add
	r.Dictionary["-"] = r.sub
	r.Dictionary["*"] = r.mul
	r.Dictionary["/"] = r.div
	r.Dictionary["MOD"] = r.mod

	// Comparison stuff
	r.Dictionary["="] = r.eq
	r.Dictionary["<"] = r.lt
	r.Dictionary[">"] = r.gt
	r.Dictionary["<="] = r.lte
	r.Dictionary[">="] = r.gte
	r.Dictionary["<>"] = r.neq

	// Conditional stuff
	r.Dictionary["IF"] = r.ifStatement
	r.Dictionary["JMP"] = r.jmp
}

func (r *Robot) jmp() error {
	skipTo, err := r.Instructions.GetNext()
	if err != nil {
		return err
	}
	r.Instructions.Ptr = int(skipTo)
	return nil
}

func (r *Robot) ifStatement() error {
	cond, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if cond.Type != T_BOOL {
		return fmt.Errorf("expected bool for IF, got %s", cond.Type)
	}
	skipTo, err := r.Instructions.GetNext()
	if err != nil {
		return err
	}

	if !cond.Value.(bool) {
		r.Instructions.Ptr = int(skipTo)
	}
	return nil
}

func (r *Robot) rot() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	c, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	r.RobotValueStack.Push(b)
	r.RobotValueStack.Push(a)
	r.RobotValueStack.Push(c)
	return nil
}

func (r *Robot) over() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	r.RobotValueStack.Push(b)
	r.RobotValueStack.Push(a)
	r.RobotValueStack.Push(b)
	return nil
}

func (r *Robot) swap() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	r.RobotValueStack.Push(a)
	r.RobotValueStack.Push(b)
	return nil
}

func (r *Robot) dup() error {
	top, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	r.RobotValueStack.Push(top)
	r.RobotValueStack.Push(top)
	return nil
}

func (r *Robot) drop() error {
	_, err := r.RobotValueStack.Pop()
	return err
}

func (r *Robot) clear() error {
	r.RobotValueStack = &stack.RobotStack[RobotValue]{}
	return nil
}

func (r *Robot) cr() error {
	fmt.Fprintln(r.Output)
	return nil
}

func (r *Robot) v() error {
	if len(*r.RobotValueStack) == 0 {
		return fmt.Errorf("stack is empty")
	}

	for _, el := range *r.RobotValueStack {
		switch el.Type {
		case T_STRING:
			fmt.Fprintf(r.Output, "%#v\n", el.Value)
		default:
			fmt.Fprintln(r.Output, el.Value)
		}
	}
	return nil
}

func (r *Robot) prn() error {
	top, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	fmt.Fprintln(r.Output, top.Value)
	return nil
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
