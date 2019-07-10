package utils

import (
	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
)

// ParseByteSize will tried to parse string to byte size.
func ParseByteSize(s string) (int64, error) {
	var v datasize.ByteSize
	err := v.UnmarshalText([]byte(s))
	if err != nil {
		log.Errorf("Expect size <%s> is invalid [%s]", s, err)
		return 0, constants.ErrorByteSizeInvalid
	}
	return int64(v), nil
}
