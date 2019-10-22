package constants

// FlowType is the type for flow
type FlowType uint8

// All available flow
const (
	FlowInvalid FlowType = iota
	FlowAtRemote
	FlowToLocal
	FlowToRemote
)

// ListType is the type for list
type ListType uint8

// All available list type
const (
	ListTypeInvalid ListType = iota
	ListTypeBucket
	ListTypeKey
)
