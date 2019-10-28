package utils

import (
	"fmt"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/pkg/fault"
)

// CalculatePartSize will calculate the object's part size.
func CalculatePartSize(size int64) (partSize int64, err error) {
	partSize = constants.DefaultPartSize

	if size > constants.MaximumObjectSize {
		return 0, fmt.Errorf("calculate part size failed: {%w}", fault.NewLocalFileTooLarge(nil, size))
	}

	for size/partSize >= int64(constants.MaximumMultipartNumber) {
		if partSize < constants.MaximumAutoMultipartSize {
			partSize = partSize << 1
			continue
		}
		// Try to adjust partSize if it is too small and account for
		// integer division truncation.
		partSize = size/int64(constants.MaximumMultipartNumber) + 1
		break
	}
	return
}
