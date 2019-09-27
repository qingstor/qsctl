package utils

import (
	"github.com/c2h5oh/datasize"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

// ParseByteSize will tried to parse string to byte size.
func ParseByteSize(s string) (int64, error) {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte(s))
	if err != nil {
		return 0, fault.NewUserInputByteSizeInvalid(err, s)
	}
	return int64(v), nil
}
