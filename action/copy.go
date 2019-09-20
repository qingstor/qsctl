package action

import (
	"bufio"
	"io"

	log "github.com/sirupsen/logrus"
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
