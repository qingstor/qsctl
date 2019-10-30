package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/c2h5oh/datasize"
	"github.com/yunify/qsctl/v2/pkg/types"
)

// ErrReadableSizeFormat will return when size format not the same with expected.
var ErrReadableSizeFormat = errors.New("readable size format invalid")

// ParseByteSize will tried to parse string to byte size.
func ParseByteSize(s string) (int64, error) {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte(s))
	if err != nil {
		return 0, types.NewErrUserInputByteSizeInvalid(err, s)
	}
	return int64(v), nil
}

// UnixReadableSize will transfer readable size string into Unix size.
// 1 KB --> 1B, 1.2 GB --> 1.2G, 103 B --> 103B
func UnixReadableSize(hrSize string) (string, error) {
	parts := strings.Split(hrSize, " ")
	if len(parts) < 2 || // no space
		!strings.ContainsRune(parts[1], 'B') || // second part does not contain 'B'
		len(parts[0]) < 1 { // no first part
		return "", types.NewErrReadableSizeFormatInvalid(ErrReadableSizeFormat, hrSize)
	}
	return fmt.Sprintf("%s%c", parts[0], parts[1][0]), nil
}
