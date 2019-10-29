package task

import (
	"fmt"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/mock"

	"github.com/yunify/qsctl/v2/pkg/types"

	"github.com/yunify/qsctl/v2/utils"
)

func TestNewStatTask(t *testing.T) {
	cases := []struct {
		input            string
		expectedTodoFunc types.TaskFunc
		expectErr        error
	}{
		{"qs://test-bucket/obj", NewObjectStatTask, nil},
		// this test case is for PseudoDir, which will return error in the near future
		// {"qs://test-bucket/obj/", NewObjectStatTask, nil},
	}

	for _, v := range cases {
		pt := NewStatTask(func(task *StatTask) {
			_, _, _, err := utils.ParseQsPath(v.input)
			if err != nil {
				t.Fatal(err)
			}
		})

		assert.Equal(t,
			fmt.Sprintf("%v", v.expectedTodoFunc),
			fmt.Sprintf("%v", pt.Todo.NextTODO()))
	}
}
func TestObjectStatTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	objectKey := uuid.New().String()
	store := mock.NewMockStorager(ctrl)

	store.EXPECT().Stat(gomock.Any(), gomock.Any()).Do(func(inputPath string, option ...*types.Pair) {
		assert.Equal(t, objectKey, inputPath)
	})

	pool := navvy.NewPool(10)

	x := &mockObjectStatTask{}
	x.SetDestinationPath(objectKey)
	x.SetPool(pool)
	x.SetDestinationStorage(store)

	task := NewObjectStatTask(x)
	task.Run()
	pool.Wait()
}
