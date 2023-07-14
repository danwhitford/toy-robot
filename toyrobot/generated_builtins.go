package toyrobot

import "fmt"

func (r *Robot) mul() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_INT, Value: b.Value.(int) * a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) add() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_INT, Value: b.Value.(int) + a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) sub() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_INT, Value: b.Value.(int) - a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) div() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_INT, Value: b.Value.(int) / a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) mod() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_INT, Value: b.Value.(int) % a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) eq() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_BOOL, Value: b.Value.(int) == a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) neq() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_BOOL, Value: b.Value.(int) != a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) lt() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_BOOL, Value: b.Value.(int) < a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) gt() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_BOOL, Value: b.Value.(int) > a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) lte() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_BOOL, Value: b.Value.(int) <= a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}

func (r *Robot) gte() error {
	a, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	b, err := r.RobotValueStack.Pop()
	if err != nil {
		return err
	}
	if a.Type != b.Type {
		return fmt.Errorf("types do not match")
	}
	switch a.Type {
	case T_INT:
		r.RobotValueStack.Push(RobotValue{Type: T_BOOL, Value: b.Value.(int) >= a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}
