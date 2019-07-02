package action

import (
	"io"

	"github.com/yunify/qsctl/constants"
)

// Copy will handle all copy actions.
func Copy(src, dest string) (err error) {
	flow, err := ParseDirection(src, dest)
	if err != nil {
		panic(err)
	}

	switch flow {
	case constants.DirectionLocalToRemote:
		r, err := ParseFilePathForRead(src)
		if err != nil {
			panic(err)
		}

		objectKey, err := ParseQsPathForWrite(dest)
		if err != nil {
			panic(err)
		}

		switch r.(type) {
		case io.ReadSeeker:
			err = CopySeekableFileToRemote(r, objectKey)
		default:
			err = CopyNotSeekableFileToRemote(r, objectKey)
		}
		if err != nil {
			panic(err)
		}
		return nil

	case constants.DirectionRemoteToLocal:
		panic("invalid flow")

	default:
		panic("invalid flow")
	}
}

// CopySeekableFileToRemote will copy a seekable file to remote.
func CopySeekableFileToRemote(r io.Reader, objectKey string) (err error) {
	return nil
}

// CopyNotSeekableFileToRemote will copy a not seekable file to remote.
func CopyNotSeekableFileToRemote(r io.Reader, objectKey string) (err error) {
	return nil
}
