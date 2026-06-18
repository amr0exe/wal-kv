package types

// OpType refers to operation performed
// It can be set|get
// 1|2 numbers are easier format to represent action
type OpType uint8

const (
	OpSet OpType = 1
	OpDel OpType = 2
)

// Node Role
// Role will define how they would act in distributed network
type Role uint8

const (
	Primary Role = iota
	Replica
)

// Snapshot structure
type SnapshotRecord struct {
	Seq   uint32
	Op    OpType
	Key   string
	Value string
}
