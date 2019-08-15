package action

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Xuanwo/navvy"
	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
)

type CopyTestSuite struct {
	suite.Suite
}

func (suite CopyTestSuite) SetupTest() {
	log.SetLevel(log.DebugLevel)

	contexts.Storage = storage.NewMockObjectStorage()
	contexts.Pool, _ = navvy.NewPool(10)
}

func (suite CopyTestSuite) TestCopyLargeFile() {
	f, err := ioutil.TempFile("", "qsctl_copy_*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	_, err = f.Seek(int64(datasize.GB), io.SeekStart)
	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte("x"))
	if err != nil {
		panic(err)
	}

	input := CopyHandler{
		FilePath: f.Name(),
	}
	input.WithObjectKey(storage.MockTBObject).RunCopyLargeFileTask()
	time.Sleep(time.Second)
	contexts.Pool.Wait()
}

func TestCopyTestSuite(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
