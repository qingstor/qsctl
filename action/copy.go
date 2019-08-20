package action

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"io"

	"github.com/Xuanwo/navvy"
	"github.com/pengsrc/go-shared/buffer"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/task/utils"
)

// CopyHandler is all params for Copy func
type CopyHandler struct {
	// Bench is whether enable benchmark
	Bench bool `json:"bench"`
	// Dest is the destination path
	Dest string `json:"dest"`
	// ExpectSize is the expect size for uploading file from stdin
	ExpectSize int64 `json:"expect_size"`
	// MaximumMemoryContent is the maximum content loaded in memory
	MaximumMemoryContent int64 `json:"maximum_memory_content"`
	// ObjectKey is the remote object key
	ObjectKey string `json:"object_key"`
	FilePath  string
	// Reader is the stream for upload
	Reader io.Reader `json:"reader"`
	// Src is the source path
	Src string `json:"src"`
	// Writer is the stream for download
	Writer io.Writer `json:"writer"`
	// Zone specifies the zone for copy action
	Zone string `json:"zone"`
}

// WithBench sets the Bench field with given bool value
func (ch *CopyHandler) WithBench(b bool) *CopyHandler {
	ch.Bench = b
	return ch
}

// WithDest sets the Dest field with given path
func (ch *CopyHandler) WithDest(path string) *CopyHandler {
	ch.Dest = path
	return ch
}

// WithExpectSize sets the ExpectSize field with given size
func (ch *CopyHandler) WithExpectSize(size int64) *CopyHandler {
	ch.ExpectSize = size
	return ch
}

// WithMaximumMemory sets the MaximumMemoryContent field with given size
func (ch *CopyHandler) WithMaximumMemory(size int64) *CopyHandler {
	ch.MaximumMemoryContent = size
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

// WithSrc sets the Src field with given path
func (ch *CopyHandler) WithSrc(path string) *CopyHandler {
	ch.Src = path
	return ch
}

// WithWriter sets the Writer field with given writer
func (ch *CopyHandler) WithWriter(w io.Writer) *CopyHandler {
	ch.Writer = w
	return ch
}

// WithZone sets the Zone field with given zone
func (ch *CopyHandler) WithZone(z string) *CopyHandler {
	ch.Zone = z
	return ch
}

// Copy will handle all copy actions.
// func (ch *CopyHandler) Copy() (err error) {
// 	flow, err := task.ParseDirection(ch.Src, ch.Dest)
// 	if err != nil {
// 		return
// 	}
//
// 	var totalSize int64
// 	if ch.Bench {
// 		f, err := os.Create("copy_profile")
// 		if err != nil {
// 			panic(err)
// 		}
// 		pprof.StartCPUProfile(f)
// 		defer pprof.StopCPUProfile()
//
// 		cur := time.Now()
// 		defer func() {
// 			elapsed := time.Since(cur)
// 			log.Debugf("Copied %s in %s, average %s/s\n",
// 				datasize.ByteSize(totalSize).HumanReadable(),
// 				elapsed,
// 				datasize.ByteSize(float64(totalSize) / elapsed.Seconds()).HumanReadable())
// 		}()
// 	}
//
// 	switch flow {
// 	case constants.DirectionLocalToRemote:
// 		r, err := task.ParseFilePathForRead(ch.Src)
// 		if err != nil {
// 			return err
// 		}
//
// 		bucketName, objectKey, err := task.ParseQsPath(ch.Dest)
// 		if err != nil {
// 			return err
// 		}
// 		if objectKey == "" {
// 			return constants.ErrorQsPathObjectKeyRequired
// 		}
// 		err = stor.SetupBucket(bucketName, ch.Zone)
// 		if err != nil {
// 			return err
// 		}
//
// 		switch x := r.(type) {
// 		case *os.File:
// 			if x == os.Stdin {
// 				totalSize, err = ch.WithObjectKey(objectKey).WithReader(r).CopyNotSeekableFileToRemote()
// 				if err != nil {
// 					return err
// 				}
// 				return nil
// 			}
// 			return constants.ErrorActionNotImplemented
// 		default:
// 			return constants.ErrorActionNotImplemented
// 		}
//
// 	case constants.DirectionRemoteToLocal:
// 		bucketName, objectKey, err := task.ParseQsPath(ch.Src)
// 		if err != nil {
// 			return err
// 		}
// 		if objectKey == "" {
// 			return constants.ErrorQsPathObjectKeyRequired
// 		}
// 		err = stor.SetupBucket(bucketName, ch.Zone)
// 		if err != nil {
// 			return err
// 		}
//
// 		w, err := task.ParseFilePathForWrite(ch.Dest)
// 		if err != nil {
// 			return err
// 		}
//
// 		switch x := w.(type) {
// 		case *os.File:
// 			if x == os.Stdout {
// 				totalSize, err = ch.WithObjectKey(objectKey).WithWriter(w).CopyObjectToNotSeekableFile()
// 				if err != nil {
// 					return err
// 				}
// 				return nil
// 			}
// 			return constants.ErrorActionNotImplemented
// 		default:
// 			return constants.ErrorActionNotImplemented
// 		}
//
// 	default:
// 		panic(constants.ErrorFlowInvalid)
// 	}
// }

// CopyNotSeekableFileToRemote will copy a not seekable file to remote.
func (ch *CopyHandler) CopyNotSeekableFileToRemote() (total int64, err error) {
	if ch.ExpectSize == 0 {
		return 0, constants.ErrorExpectSizeRequired
	}

	uploadID, err := stor.InitiateMultipartUpload(ch.ObjectKey)
	if err != nil {
		return
	}

	log.Debugf("Object <%s> uploading via upload ID <%s>", ch.ObjectKey, uploadID)

	partSize, err := utils.CalculatePartSize(ch.ExpectSize)
	if err != nil {
		return
	}

	pool, err := navvy.NewPool(utils.CalculateConcurrentWorkers(partSize, ch.MaximumMemoryContent))
	if err != nil {
		panic(err)
	}
	defer pool.Release()

	bytesPool := buffer.NewBytesPool()

	partNumber := 0
	for {
		lr := bufio.NewReader(io.LimitReader(ch.Reader, partSize))
		b := bytesPool.Get()
		n, err := io.Copy(b, lr)

		if ch.Bench {
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
		pool.Submit(navvy.TaskWrapper(func() {
			// We should free the bytes after upload.
			defer b.Free()

			md5sum := md5.Sum(b.Bytes())

			err = stor.UploadMultipart(ch.ObjectKey, uploadID, int64(n), localPartNumber, md5sum[:], bytes.NewReader(b.Bytes()))
			if err != nil {
				log.Errorf("Object <%s> part <%d> upload failed [%s]", ch.ObjectKey, localPartNumber, err)
			}
			log.Debugf("Object <%s> part <%d> uploaded", ch.ObjectKey, localPartNumber)
		}))

		partNumber++
	}

	pool.Wait()

	err = stor.CompleteMultipartUpload(ch.ObjectKey, uploadID, partNumber)
	if err != nil {
		return
	}
	log.Infof("Object <%s> upload finished", ch.ObjectKey)
	return total, nil
}

// CopyObjectToNotSeekableFile will copy an object to not seekable file.
func (ch *CopyHandler) CopyObjectToNotSeekableFile() (total int64, err error) {
	r, err := stor.GetObject(ch.ObjectKey)
	if err != nil {
		return
	}

	bw, br := bufio.NewWriter(ch.Writer), bufio.NewReader(r)
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
