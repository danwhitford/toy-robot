// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package toyrobot

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TOKEN_NUMBER-0]
	_ = x[TOKEN_DIRECTION-1]
	_ = x[TOKEN_WORD-2]
	_ = x[TOKEN_BOOL-3]
	_ = x[TOKEN_STRING-4]
}

const _TokenType_name = "TOKEN_NUMBERTOKEN_DIRECTIONTOKEN_WORDTOKEN_BOOLTOKEN_STRING"

var _TokenType_index = [...]uint8{0, 12, 27, 37, 47, 59}

func (i TokenType) String() string {
	if i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
