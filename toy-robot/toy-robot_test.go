package toyrobot

import (
	"bytes"
	"testing"
)

func TestCanPlace(t *testing.T) {
	table := []struct {
		x, y int
		f    Direction
	}{
		{0, 0, NORTH},
		{0, 0, EAST},
		{1, 2, SOUTH},
		{2, 1, WEST},
	}

	for _, tst := range table {
		robot := NewRobot()
		if robot.Placed {
			t.Error("Robot should not be placed")
		}
		robot.place(tst.x, tst.y, tst.f)
		if !robot.Placed {
			t.Error("Robot should be placed")
		}
		if robot.X != tst.x {
			t.Error("Robot X should be 0")
		}
		if robot.Y != tst.y {
			t.Error("Robot Y should be 0")
		}
		if robot.F != tst.f {
			t.Error("Robot F should be NORTH")
		}
	}
}

func TestWontPlaceOffBoard(t *testing.T) {
	table := []struct {
		x, y int
		f    Direction
	}{
		{-1, 0, NORTH},
		{0, -1, EAST},
		{5, 0, SOUTH},
		{0, 5, WEST},
	}

	for _, tst := range table {
		robot := NewRobot()
		if robot.Placed {
			t.Error("Robot should not be placed")
		}
		robot.place(tst.x, tst.y, tst.f)
		if robot.Placed {
			t.Error("Robot should not be placed")
		}
	}
}

func TestReadPlaceInstruction(t *testing.T) {
	table := []struct {
		instruction string
		x, y        int
		f           Direction
	}{
		{"PLACE 0,0,NORTH", 0, 0, NORTH},
		{"PLACE 0,0,EAST", 0, 0, EAST},
		{"PLACE 1,2,SOUTH", 1, 2, SOUTH},
		{"PLACE 2,1,WEST", 2, 1, WEST},
	}

	for _, tst := range table {
		robot := NewRobot()
		if robot.Placed {
			t.Error("Robot should not be placed")
		}
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Errorf("Error reading instruction %s: %s", tst.instruction, err)
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
		x, y        int
		f           Direction
		x2, y2      int
		f2          Direction
	}{
		{"MOVE", 0, 1, NORTH, 0, 2, NORTH},
		{"MOVE", 1, 0, EAST, 2, 0, EAST},
		{"MOVE", 1, 2, SOUTH, 1, 1, SOUTH},
		{"MOVE", 1, 0, WEST, 0, 0, WEST},
		{"MOVE", 0, 0, SOUTH, 0, 0, SOUTH},
		{"MOVE", 0, 4, NORTH, 0, 4, NORTH},
		{"MOVE", 0, 0, WEST, 0, 0, WEST},
		{"MOVE", 4, 0, EAST, 4, 0, EAST},
	}

	for _, tst := range table {
		robot := NewRobot()
		robot.place(tst.x, tst.y, tst.f)
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
		x, y        int
		f           Direction
		f2          Direction
	}{
		{"LEFT", 0, 0, NORTH, WEST},
		{"LEFT", 0, 0, EAST, NORTH},
		{"LEFT", 0, 0, SOUTH, EAST},
		{"LEFT", 0, 0, WEST, SOUTH},
		{"RIGHT", 0, 0, NORTH, EAST},
		{"RIGHT", 0, 0, EAST, SOUTH},
		{"RIGHT", 0, 0, SOUTH, WEST},
		{"RIGHT", 0, 0, WEST, NORTH},
	}

	for _, tst := range table {
		robot := NewRobot()
		robot.place(tst.x, tst.y, tst.f)
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
		x, y        int
		f           Direction
		report      string
	}{
		{"REPORT", 0, 0, NORTH, "0,0,NORTH\n"},
		{"REPORT", 0, 0, EAST, "0,0,EAST\n"},
		{"REPORT", 0, 0, SOUTH, "0,0,SOUTH\n"},
		{"REPORT", 0, 0, WEST, "0,0,WEST\n"},
	}

	for _, tst := range table {
		var buffer bytes.Buffer
		robot := NewRobot()
		robot.Output = &buffer
		robot.place(tst.x, tst.y, tst.f)
		err := robot.ReadInstruction(tst.instruction)
		if err != nil {
			t.Errorf("Error reading instruction %s: %s", tst.instruction, err)
		}
		if !robot.Placed {
			t.Error("Robot should be placed")
		}
		reported := buffer.String()
		if reported != tst.report {
			t.Errorf("Robot report should be %s but was %s", tst.report, reported)
		}
	}
}
