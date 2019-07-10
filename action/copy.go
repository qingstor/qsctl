package action

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/c2h5oh/datasize"
	"github.com/panjf2000/ants"
	"github.com/pengsrc/go-shared/buffer"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
)

// Copy will handle all copy actions.
func Copy(src, dest string) (err error) {
	flow, err := ParseDirection(src, dest)
	if err != nil {
		return
	}

	var totalSize int64
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
			log.Debugf("Copied %s in %s, average %s/s\n",
				datasize.ByteSize(totalSize).HumanReadable(),
				elapsed,
				datasize.ByteSize(float64(totalSize)/elapsed.Seconds()).HumanReadable())
		}()
	}

	switch flow {
	case constants.DirectionLocalToRemote:
		r, err := ParseFilePathForRead(src)
		if err != nil {
			return err
		}

		bucketName, objectKey, err := ParseQsPath(dest)
		if err != nil {
			return err
		}
		if objectKey == "" {
			return constants.ErrorQsPathObjectKeyRequired
		}
		err = contexts.Storage.SetupBucket(bucketName, "")
		if err != nil {
			return err
		}

		switch x := r.(type) {
		case *os.File:
			if x == os.Stdin {
				totalSize, err = CopyNotSeekableFileToRemote(r, objectKey)
				if err != nil {
					return err
				}
				return nil
			}
			return constants.ErrorActionNotImplemented
		default:
			return constants.ErrorActionNotImplemented
		}

	case constants.DirectionRemoteToLocal:
		bucketName, objectKey, err := ParseQsPath(dest)
		if err != nil {
			return err
		}
		if objectKey == "" {
			return constants.ErrorQsPathObjectKeyRequired
		}
		err = contexts.Storage.SetupBucket(bucketName, "")
		if err != nil {
			return err
		}

		w, err := ParseFilePathForWrite(dest)
		if err != nil {
			return err
		}

		switch x := w.(type) {
		case *os.File:
			if x == os.Stdout {
				totalSize, err = CopyObjectToNotSeekableFile(w, objectKey)
				if err != nil {
					return err
				}
				return nil
			}
			return constants.ErrorActionNotImplemented
		default:
			return constants.ErrorActionNotImplemented
		}

	default:
		panic(constants.ErrorFlowInvalid)
	}
}

// CopyNotSeekableFileToRemote will copy a not seekable file to remote.
func CopyNotSeekableFileToRemote(r io.Reader, objectKey string) (total int64, err error) {
	if contexts.ExpectSize == 0 {
		return 0, constants.ErrorExpectSizeRequired
	}

	uploadID, err := contexts.Storage.InitiateMultipartUpload(objectKey)
	if err != nil {
		return
	}

	log.Debugf("Object <%s> uploading via upload ID <%s>", objectKey, uploadID)

	partSize, err := CalculatePartSize(contexts.ExpectSize)
	if err != nil {
		return
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
			total += int64(n)
		}

		if n == 0 {
			break
		}
		if err != nil {
			log.Errorf("Read failed [%v]", err)
			return 0, err
		}

		localPartNumber := partNumber
		wg.Add(1)
		err = pool.Submit(func() {
			defer wg.Done()

			// We should free the bytes after upload.
			defer b.Free()

			err = contexts.Storage.UploadMultipart(objectKey, uploadID, int64(n), localPartNumber, md5.New().Sum(b.Bytes()), bytes.NewReader(b.Bytes()))
			if err != nil {
				log.Errorf("Object <%s> part <%d> upload failed [%s]", objectKey, localPartNumber, err)
			}
			log.Debugf("Object <%s> part <%d> uploaded", objectKey, localPartNumber)
		})
		if err != nil {
			panic(err)
		}

		partNumber++
	}

	wg.Wait()

	err = contexts.Storage.CompleteMultipartUpload(objectKey, uploadID, partNumber)
	if err != nil {
		return
	}
	log.Infof("Object <%s> upload finished", objectKey)
	return total, nil
}

// CopyObjectToNotSeekableFile will copy an object to not seekable file.
func CopyObjectToNotSeekableFile(w io.Writer, objectKey string) (total int64, err error) {
	r, err := contexts.Storage.GetObject(objectKey)
	if err != nil {
		return
	}

	bw, br := bufio.NewWriter(w), bufio.NewReader(r)
	total, err = io.Copy(bw, br)
	if err != nil {
		log.Errorf("Copy failed [%v]", err)
		return 0, err
	}
	err = bw.Flush()
	if err != nil {
		log.Errorf("Buffer flush failed [%v]", err)
		return 0, err
	}
	return
}
