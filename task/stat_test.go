package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/storage"
	"github.com/yunify/qsctl/v2/utils"
)

func TestNewStatTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TodoFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", NewStatObjectTask, nil},
		{"qs://test-bucket/obj/", NewStatObjectTask, nil},
	}

	for _, v := range cases {
		pt := NewStatTask(func(task *StatTask) {
			_, _, _, err := utils.ParseKey(v.input)
			if err != nil {
				t.Fatal(err)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}

func TestStatObjectTask_Run(t *testing.T) {
	bucketName, objectKey, zone := uuid.New().String(), storage.MockMBObject, "t1"
	store := storage.NewMockObjectStorage()
	err := store.SetupBucket(bucketName, zone)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := navvy.NewPool(10)
	if err != nil {
		t.Fatal(err)
	}

	x := &mockStatObjectTask{}
	x.SetKey(objectKey)
	x.SetPool(pool)
	x.SetStorage(store)
	x.SetObjectMeta(&storage.ObjectMeta{})

	task := NewStatObjectTask(x)
	task.Run()
	pool.Wait()

	cases := []struct {
		input          string
		expectedLength int64
		expectErr      error
	}{
		{storage.Mock0BObject, int64(0), nil},
		{storage.MockMBObject, int64(1024 * 1024), nil},
		{storage.MockGBObject, int64(1024 * 1024 * 1024), nil},
		{storage.MockTBObject, int64(1024 * 1024 * 1024 * 1024), nil},
	}

	for _, v := range cases {
		om, err := store.HeadObject(v.input)
		assert.Equal(t, v.expectErr, err)
		assert.Equal(t, v.expectedLength, om.ContentLength)
	}

}
