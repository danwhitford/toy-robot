package toyrobot

//go:generate stringer -type=RobotType
type RobotType byte

const (
	T_INT RobotType = iota
	T_DIRECTION
	T_BOOL
	T_STRING
)

type RobotValue struct {
	Type  RobotType
	Value any
}

//go:generate stringer -type=Direction
type Direction byte

const (
	NORTH Direction = iota
	EAST
	SOUTH
	WEST
)
