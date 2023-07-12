package toyrobot

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed programs/*.bot
var programs embed.FS

func TestWontPlaceOffBoard(t *testing.T) {
	table := []struct {
		input string
	}{
		{"7 0 NORTH PLACE"},
		{"0 7 EAST PLACE"},
		{"5 0 SOUTH PLACE"},
		{"0 5 WEST PLACE"},
	}

	for _, tst := range table {
		robot := NewRobot()
		if robot.Placed {
			t.Errorf("Robot should not be placed but is at %d %d %s", robot.X, robot.Y, robot.F)
		}
		err := robot.ReadInstruction(tst.input)
		if err != nil {
			t.Fatalf("error for %s: %s", tst.input, err)
		}
		if robot.Placed {
			t.Fatalf("Robot should not be placed but is at %d %d %s", robot.X, robot.Y, robot.F)
		}
	}
}

func TestReadPlaceInstruction(t *testing.T) {
	table := []struct {
		instruction string
		x, y        int
		f           Direction
	}{
		{"0 0 NORTH PLACE", 0, 0, NORTH},
		{"0 0 EAST PLACE", 0, 0, EAST},
		{"1 2 SOUTH PLACE", 1, 2, SOUTH},
		{"2 1 WEST PLACE", 2, 1, WEST},
	}

	for _, tst := range table {
		robot := NewRobot()
		if robot.Placed {
			t.Error("Robot should not be placed")
		}
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Fatalf("Error reading instruction %s: %s", tst.instruction, err)
		}
		if !robot.Placed {
			t.Error("Robot should be placed")
		}
		if robot.X != tst.x {
			t.Errorf("Robot X should be %d", tst.x)
		}
		if robot.Y != tst.y {
			t.Errorf("Robot Y should be %d", tst.y)
		}
		if robot.F != tst.f {
			t.Errorf("Robot F should be %d", tst.f)
		}
	}
}

func TestReadMoveInstruction(t *testing.T) {
	table := []struct {
		instruction string
		x2, y2      int
		f2          Direction
	}{
		{"1 0 EAST PLACE MOVE", 2, 0, EAST},
		{"1 2 SOUTH PLACE MOVE", 1, 1, SOUTH},
		{"1 0 WEST PLACE MOVE", 0, 0, WEST},
		{"0 0 SOUTH PLACE MOVE", 0, 0, SOUTH},
		{"0 4 NORTH PLACE MOVE", 0, 4, NORTH},
		{"0 0 WEST PLACE MOVE", 0, 0, WEST},
		{"4 0 EAST PLACE MOVE", 4, 0, EAST},
	}

	for _, tst := range table {
		robot := NewRobot()
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Errorf("Error reading instruction %s: %s", tst.instruction, err)
		}
		if !robot.Placed {
			t.Error("Robot should be placed")
		}
		if robot.X != tst.x2 {
			t.Errorf("Robot X should be %d but was %d", tst.x2, robot.X)
		}
		if robot.Y != tst.y2 {
			t.Errorf("Robot Y should be %d but was %d", tst.y2, robot.Y)
		}
		if robot.F != tst.f2 {
			t.Errorf("Robot F should be %d but was %d", tst.f2, robot.F)
		}
	}
}

func TestLeftRightInstructions(t *testing.T) {
	table := []struct {
		instruction string
		f2          Direction
	}{
		{"0 0 NORTH PLACE LEFT", WEST},
		{"0 0 EAST PLACE LEFT", NORTH},
		{"0 0 SOUTH PLACE LEFT", EAST},
		{"0 0 WEST PLACE LEFT", SOUTH},
		{"0 0 NORTH PLACE RIGHT", EAST},
		{"0 0 EAST PLACE RIGHT", SOUTH},
		{"0 0 SOUTH PLACE RIGHT", WEST},
		{"0 0 WEST PLACE RIGHT", NORTH},
	}

	for _, tst := range table {
		robot := NewRobot()
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Errorf("Error reading instruction %s: %s", tst.instruction, err)
		}
		if !robot.Placed {
			t.Error("Robot should be placed")
		}
		if robot.F != tst.f2 {
			t.Errorf("Robot F should be %d but was %d", tst.f2, robot.F)
		}
	}
}

func TestReportInstruction(t *testing.T) {
	table := []struct {
		instruction string
		report      string
	}{
		{"0 0 NORTH PLACE REPORT", "0,0,NORTH\n"},
		{"0 0 EAST PLACE REPORT", "0,0,EAST\n"},
		{"0 0 SOUTH PLACE REPORT", "0,0,SOUTH\n"},
		{"0 0 WEST PLACE REPORT", "0,0,WEST\n"},
		{"REPORT", "Robot not placed\n"},
	}

	for _, tst := range table {
		var buffer bytes.Buffer
		robot := NewRobot()
		robot.Output = &buffer
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Fatalf("Error reading instruction %s: %s", tst.instruction, err)
		}
		reported := buffer.String()
		if reported != tst.report {
			t.Errorf("Robot report should be %s but was %s", tst.report, reported)
		}
	}
}

func TestManyInstructionsOnOneLine(t *testing.T) {
	table := []struct {
		instruction    string
		expectedReport string
	}{
		{"0 0 NORTH PLACE MOVE LEFT MOVE REPORT", "0,1,WEST\n"},
		{"1 2 EAST PLACE MOVE MOVE LEFT MOVE REPORT", "3,3,NORTH\n"},
		{"0 0 NORTH PLACE LEFT LEFT LEFT LEFT REPORT", "0,0,NORTH\n"},
		{"0 0 NORTH PLACE RIGHT RIGHT RIGHT RIGHT REPORT", "0,0,NORTH\n"},
		{"1 2 EAST PLACE MOVE MOVE LEFT MOVE REPORT", "3,3,NORTH\n"},
		{"0 0 NORTH PLACE MOVE LEFT MOVE REPORT", "0,1,WEST\n"},
		{"0 0 NORTH PLACE MOVE LEFT MOVE\nREPORT # SOME COMPLETE NONSENSE", "0,1,WEST\n"},
	}

	for i, tst := range table {
		var buffer bytes.Buffer
		robot := NewRobot()
		robot.Output = &buffer
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Errorf("Error reading instruction %s: %s", tst.instruction, err)
		}
		reported := buffer.String()
		if reported != tst.expectedReport {
			tst.expectedReport = strings.Replace(tst.expectedReport, "\n", "\\n", -1)
			reported = strings.Replace(reported, "\n", "\\n", -1)
			t.Errorf("Failed %dth test. Robot report should be '%s' but was '%s'", i, tst.expectedReport, reported)
		}
	}
}

func TestWholePrograms(t *testing.T) {
	testEnts, err := programs.ReadDir("programs")
	if err != nil {
		t.Fatalf("Error reading test programs: %s", err)
	}

	var table []struct {
		program        string
		expectedOutput string
		fname          string
	}
	for _, testEnt := range testEnts {
		name := fmt.Sprintf("programs/%s", testEnt.Name())
		contentBytes, err := programs.ReadFile(name)
		if err != nil {
			t.Errorf("Error reading test program %s: %s", testEnt.Name(), err)
		}
		content := string(contentBytes)
		parts := strings.Split(content, "### OUTPUT ###\n")
		program := parts[0]
		var outputSB strings.Builder
		for _, line := range strings.Split(parts[1], "\n") {
			line = strings.TrimPrefix(line, "# ")
			outputSB.WriteString(line)
			outputSB.WriteString("\n")
		}
		table = append(table, struct {
			program        string
			expectedOutput string
			fname          string
		}{
			program,
			outputSB.String(),
			testEnt.Name(),
		})
	}

	for _, tst := range table {
		var buffer bytes.Buffer
		robot := NewRobot()
		robot.Output = &buffer
		err := robot.ReadInstruction(tst.program)
		if err != nil {
			t.Fatalf("Error reading program %s: %s", tst.program, err)
		}
		got := buffer.String()
		if diff := cmp.Diff(tst.expectedOutput, got); diff != "" {
			t.Errorf("Program output mismatch for '%s' (-want +got):\n%s", tst.fname, diff)
		}
	}
}
