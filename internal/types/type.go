package types

// OpType refers to operation performed
// It can be set|get
// 1|2 numbers are easier format to represent action
type OpType uint8

const (
	OpSet OpType = 1
	OpDel OpType = 2
)
