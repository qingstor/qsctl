package task

import (
	"os"
	"sync"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

func TestCopyLargeFileTask_Run(t *testing.T) {
	bucketName := uuid.New().String()
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, "")
	if err != nil {
		t.Fatal(err)
	}
	key := uuid.New().String()

	name, size, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopyLargeFileTask{}
	x.SetPool(pool)
	x.SetPath(name)
	x.SetKey(key)
	x.SetStorage(store)
	x.SetTotalSize(size)

	task := NewCopyLargeFileTask(x)
	task.Run()
	pool.Wait()

	om, err := store.HeadObject(key)
	assert.NoError(t, err)
	assert.Equal(t, size, om.ContentLength)
}

func TestCopyPartialFileTask_Run(t *testing.T) {
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

	name, size, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopyPartialFileTask{}
	x.SetPool(pool)
	x.SetPath(name)
	x.SetKey(key)
	x.SetStorage(store)
	x.SetUploadID(uploadID)
	x.SetPartSize(64 * 1024 * 1024)
	x.SetTotalSize(size)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	x.SetWaitGroup(wg)

	currentPartNumber := int32(0)
	x.SetCurrentPartNumber(&currentPartNumber)

	currentOffset := int64(0)
	x.SetCurrentOffset(&currentOffset)

	task := NewCopyPartialFileTask(x)
	task.Run()
	pool.Wait()

	multipart, ok := store.Multipart[key]
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(multipart.Parts))
	assert.Equal(t, size, multipart.Length)
}

func TestCopySmallFileTask_Run(t *testing.T) {
	bucketName := uuid.New().String()
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, "")
	if err != nil {
		t.Fatal(err)
	}
	key := uuid.New().String()

	name, size, _ := utils.GenerateTestFile()
	defer os.Remove(name)

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockCopySmallFileTask{}
	x.SetPool(pool)
	x.SetPath(name)
	x.SetKey(key)
	x.SetStorage(store)
	x.SetTotalSize(size)

	task := NewCopySmallFileTask(x)
	task.Run()
	pool.Wait()

	om, err := store.HeadObject(key)
	assert.NoError(t, err)
	assert.Equal(t, size, om.ContentLength)
}
