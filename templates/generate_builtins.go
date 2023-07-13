package main

import (
	"fmt"
	"os"
	"text/template"
)

const robotTemplate = `package toyrobot

import "fmt"

{{ range . }}
func (r *Robot) {{ .FunctionName }}() error {
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
		r.RobotValueStack.Push(RobotValue{Type: T_INT, Value: b.Value.(int) {{ .FunctionOp }} a.Value.(int)})
	default:
		return fmt.Errorf("unsupported type")
	}
	return nil
}
{{ end }}
`

func main() {
	datas := []struct {
		FunctionName string
		FunctionOp   string
	}{
		{"mul", "*"},
		{"add", "+"},
		{"sub", "-"},
		{"div", "/"},
		{"eq", "=="},
		{"neq", "!="},
		{"lt", "<"},
		{"gt", ">"},
		{"lte", "<="},
		{"gte", ">="},
	}

	tmpl, err := template.New("robotMul").Parse(robotTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	f, err := os.Create("./generated_builtins.go")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()
	err = tmpl.ExecuteTemplate(f, "robotMul", datas)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}
