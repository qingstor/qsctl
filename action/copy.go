package action

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/helper"
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

		switch x := r.(type) {
		case *os.File:
			if x.Name() == "/dev/stdin" {
				err = CopyNotSeekableFileToRemote(r, objectKey)
				return err
			}
			err = CopySeekableFileToRemote(r, objectKey)
			return err
		case io.ReadSeeker:
			fmt.Printf("Start CopySeekableFileToRemote")
			err = CopySeekableFileToRemote(r, objectKey)
			if err != nil {
				panic(err)
			}
		default:
			fmt.Printf("Start CopyNotSeekableFileToRemote")
			err = CopyNotSeekableFileToRemote(r, objectKey)
			if err != nil {
				panic(err)
			}
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
	uploadID, err := helper.InitiateMultipartUpload(objectKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload ID is %s.\n", uploadID)

	partNumber := 0

	for {
		lr := io.LimitReader(r, constants.DefaultPartSize)
		b, err := ioutil.ReadAll(lr)
		l := len(b)
		if l == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Printf("Read %d bytes.\n", l)
		err = helper.UploadMultipart(objectKey, uploadID, int64(l), partNumber, md5.Sum(b), bytes.NewReader(b))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Part %d uploaded.\n", partNumber)

		partNumber++
	}

	err = helper.CompleteMultipartUpload(objectKey, uploadID, partNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload ID %s for %s finished.\n", uploadID, objectKey)
	return nil
}
