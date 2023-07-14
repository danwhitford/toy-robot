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
