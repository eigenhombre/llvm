// Code generated by "stringer -linecomment -type SelectionKind"; DO NOT EDIT.

package enum

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SelectionKindAny-0]
	_ = x[SelectionKindExactMatch-1]
	_ = x[SelectionKindLargest-2]
	_ = x[SelectionKindNoDeduplicate-3]
	_ = x[SelectionKindSameSize-4]
}

const _SelectionKind_name = "anyexactmatchlargestnodeduplicatesamesize"

var _SelectionKind_index = [...]uint8{0, 3, 13, 20, 33, 41}

func (i SelectionKind) String() string {
	if i >= SelectionKind(len(_SelectionKind_index)-1) {
		return "SelectionKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SelectionKind_name[_SelectionKind_index[i]:_SelectionKind_index[i+1]]
}
