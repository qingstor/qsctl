package task

import (
	"bytes"
	"sync"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
	taskUtil "github.com/yunify/qsctl/v2/task/utils"
)

func TestCopyStreamTask_Run(t *testing.T) {
	bucketName := uuid.New().String()
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, "")
	if err != nil {
		t.Fatal(err)
	}
	key := uuid.New().String()

	buf, size, _ := taskUtil.GenerateTestStream()

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := NewCopyTask(func(task *CopyTask) {
		task.SetStorage(store)
		task.SetKey(key)
		task.SetPath("-")
		task.SetPool(pool)
		task.SetKeyType(constants.KeyTypeObject)
		task.SetPathType(constants.PathTypeStream)
		task.SetFlowType(constants.FlowToRemote)
		task.SetStream(buf)
	})

	task := NewCopyStreamTask(x)
	task.Run()
	pool.Wait()

	object, ok := store.Meta[key]
	assert.Equal(t, true, ok)
	assert.Equal(t, size, object.ContentLength)
}

func TestCopyPartialStreamTask_Run(t *testing.T) {
	bucketName := uuid.New().String()
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, "")
	if err != nil {
		t.Fatal(err)
	}
	key := uuid.New().String()

	uploadID, err := store.InitiateMultipartUpload(key)
	if err != nil {
		t.Fatal(err)
	}

	buf, size, _ := taskUtil.GenerateTestStream()

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopyPartialStreamTask{}
	x.SetPool(pool)
	x.SetStream(buf)
	x.SetKey(key)
	x.SetStorage(store)
	x.SetUploadID(uploadID)
	x.SetPartSize(64 * 1024 * 1024)
	x.SetBytesPool(&sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, x.GetPartSize()))
		},
	})

	wg := &sync.WaitGroup{}
	wg.Add(1)
	x.SetWaitGroup(wg)

	currentPartNumber := int32(0)
	x.SetCurrentPartNumber(&currentPartNumber)

	currentOffset := int64(0)
	x.SetCurrentOffset(&currentOffset)

	task := NewCopyPartialStreamTask(x)
	task.Run()
	pool.Wait()

	multipart, ok := store.Multipart[key]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(multipart.Parts))
	assert.Equal(t, size, multipart.Length)
}
