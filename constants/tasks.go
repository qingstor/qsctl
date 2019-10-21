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

// PathType is the type for path
type PathType uint8

// All available path type
const (
	PathTypeInvalid PathType = iota
	PathTypeFile
	PathTypeStream
	PathTypeLocalDir
)

// KeyType is the type for key
type KeyType uint8

// All available key type
const (
	KeyTypeInvalid KeyType = iota
	KeyTypeBucket
	KeyTypeObject
	KeyTypePseudoDir
)

// ListType is the type for list
type ListType uint8

// All available list type
const (
	ListTypeInvalid ListType = iota
	ListTypeBucket
	ListTypeKey
)
