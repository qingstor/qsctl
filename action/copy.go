package action

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/panjf2000/ants"
	"github.com/pengsrc/go-shared/buffer"

	"github.com/yunify/qsctl/constants"
	"github.com/yunify/qsctl/contexts"
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
	totalSize := int64(0)
	if contexts.Bench {
		f, err := os.Create("profile")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		cur := time.Now()
		defer func() {
			elapsed := time.Since(cur)
			fmt.Printf("Copied %s in %s, avgerage %s/s\n",
				datasize.ByteSize(totalSize).HumanReadable(),
				elapsed,
				datasize.ByteSize(float64(totalSize)/elapsed.Seconds()).HumanReadable())
		}()
	}
	if contexts.ExpectSize == 0 {
		panic("invalid expect size")
	}

	uploadID, err := helper.InitiateMultipartUpload(objectKey)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload ID is %s.\n", uploadID)

	partSize, err := CalculatePartSize(contexts.ExpectSize)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	pool, err := ants.NewPool(CalculateConcurrentWorkers(partSize))
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	bytesPool := buffer.NewBytesPool()

	partNumber := 0

	for {
		lr := bufio.NewReader(io.LimitReader(r, partSize))
		b := bytesPool.Get()
		n, err := io.Copy(b, lr)
		if contexts.Bench {
			totalSize += int64(n)
		}
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Printf("Read %d bytes.\n", n)

		localPartNumber := partNumber
		err = pool.Submit(func() {
			wg.Add(1)
			defer wg.Done()

			// We should free the bytes after upload.
			defer b.Free()

			err = helper.UploadMultipart(objectKey, uploadID, int64(n), localPartNumber, md5.Sum(b.Bytes()), bytes.NewReader(b.Bytes()))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Part %d uploaded.\n", localPartNumber)
		})
		if err != nil {
			panic(err)
		}

		partNumber++
	}

	wg.Wait()

	err = helper.CompleteMultipartUpload(objectKey, uploadID, partNumber)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload ID %s for %s finished.\n", uploadID, objectKey)
	return nil
}
