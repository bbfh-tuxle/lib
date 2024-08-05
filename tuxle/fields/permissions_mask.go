package fields

import "strings"

type permMask string

const (
	CAN_NOTHING  permMask = "----"
	CAN_READ     permMask = "r---"
	CAN_INTERACT permMask = "ri--"
	CAN_MODIFY   permMask = "rim-"
	CAN_ALL      permMask = "rim*"
)

func (perm permMask) Mask() string {
	return strings.ReplaceAll(string(perm), "-", "")
}

func (perm permMask) Can(mask permMask) bool {
	return strings.Contains(string(perm), mask.Mask())
}
