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

// CopyHandler is all params for Copy func
type CopyHandler struct {
	*FlagHandler
	// Src is the source path
	Src string `json:"src"`
	// Dest is the destination path
	Dest string `json:"dest"`
	// ObjectKey is the remote object key
	ObjectKey string `json:"object_key"`
	// Reader is the stream for upload
	Reader io.Reader `json:"reader"`
	// Writer is the stream for download
	Writer io.Writer `json:"writer"`
}

// WithBench rewrite the WithBench method
func (ch *CopyHandler) WithBench(b bool) *CopyHandler {
	ch.FlagHandler = ch.FlagHandler.WithBench(b)
	return ch
}

// WithExpectSize rewrite the WithExpectSize method
func (ch *CopyHandler) WithExpectSize(size int64) *CopyHandler {
	ch.FlagHandler = ch.FlagHandler.WithExpectSize(size)
	return ch
}

// WithMaximumMemory rewrite the WithMaximumMemory method
func (ch *CopyHandler) WithMaximumMemory(size int64) *CopyHandler {
	ch.FlagHandler = ch.FlagHandler.WithMaximumMemory(size)
	return ch
}

// WithZone rewrite the WithZone method
func (ch *CopyHandler) WithZone(z string) *CopyHandler {
	ch.FlagHandler = ch.FlagHandler.WithZone(z)
	return ch
}

// WithSrc sets the Src field with given path
func (ch *CopyHandler) WithSrc(path string) *CopyHandler {
	ch.Src = path
	return ch
}

// WithDest sets the Dest field with given path
func (ch *CopyHandler) WithDest(path string) *CopyHandler {
	ch.Dest = path
	return ch
}

// WithObjectKey sets the ObjectKey field with given key
func (ch *CopyHandler) WithObjectKey(key string) *CopyHandler {
	ch.ObjectKey = key
	return ch
}

// WithReader sets the Reader field with given reader
func (ch *CopyHandler) WithReader(r io.Reader) *CopyHandler {
	ch.Reader = r
	return ch
}

// WithWriter sets the Writer field with given writer
func (ch *CopyHandler) WithWriter(w io.Writer) *CopyHandler {
	ch.Writer = w
	return ch
}

// Copy will handle all copy actions.
func (ch *CopyHandler) Copy() (err error) {
	// Get params from handler
	bench := ch.Bench
	zone := ch.Zone
	src := ch.Src
	dest := ch.Dest

	flow, err := ParseDirection(src, dest)
	if err != nil {
		return
	}

	var totalSize int64
	if bench {
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
		err = contexts.Storage.SetupBucket(bucketName, zone)
		if err != nil {
			return err
		}

		switch x := r.(type) {
		case *os.File:
			if x == os.Stdin {
				totalSize, err = ch.WithObjectKey(objectKey).WithReader(r).CopyNotSeekableFileToRemote()
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
		bucketName, objectKey, err := ParseQsPath(src)
		if err != nil {
			return err
		}
		if objectKey == "" {
			return constants.ErrorQsPathObjectKeyRequired
		}
		err = contexts.Storage.SetupBucket(bucketName, zone)
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
				totalSize, err = ch.WithObjectKey(objectKey).WithWriter(w).CopyObjectToNotSeekableFile()
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
func (ch *CopyHandler) CopyNotSeekableFileToRemote() (total int64, err error) {
	// Get params from handler
	bench := ch.Bench
	expectSize := ch.ExpectSize
	maximumMemory := ch.MaximumMemoryContent
	objectKey := ch.ObjectKey
	r := ch.Reader

	if expectSize == 0 {
		return 0, constants.ErrorExpectSizeRequired
	}

	uploadID, err := contexts.Storage.InitiateMultipartUpload(objectKey)
	if err != nil {
		return
	}

	log.Debugf("Object <%s> uploading via upload ID <%s>", objectKey, uploadID)

	partSize, err := CalculatePartSize(expectSize)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	pool, err := ants.NewPool(CalculateConcurrentWorkers(partSize, maximumMemory))
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

		if bench {
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

			md5sum := md5.Sum(b.Bytes())

			err = contexts.Storage.UploadMultipart(objectKey, uploadID, int64(n), localPartNumber, md5sum[:], bytes.NewReader(b.Bytes()))
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
func (ch *CopyHandler) CopyObjectToNotSeekableFile() (total int64, err error) {
	// Get params from handler
	objectKey := ch.ObjectKey
	w := ch.Writer

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
