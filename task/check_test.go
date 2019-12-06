package task

import (
	"testing"
	"time"

	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/pkg/fault"
	"github.com/qingstor/qsctl/v2/pkg/mock"
)

func TestBetweenStorageCheckTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name         string
		expectObject *typ.Object
		expectErr    error
	}{
		{
			"normal",
			&typ.Object{},
			nil,
		},
		{
			"dst object not exist",
			nil,
			typ.ErrObjectNotExist,
		},
		{
			"error",
			nil,
			typ.ErrUnhandledError,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			srcStore := mock.NewMockStorager(ctrl)
			dstStore := mock.NewMockStorager(ctrl)
			srcPath := uuid.New().String()
			dstPath := uuid.New().String()

			task := BetweenStorageCheckTask{}
			task.SetSourceStorage(srcStore)
			task.SetDestinationStorage(dstStore)
			task.SetSourcePath(srcPath)
			task.SetDestinationPath(dstPath)
			task.SetFault(fault.New())

			srcStore.EXPECT().Stat(gomock.Any()).DoAndReturn(func(path string, pairs ...*typ.Pair) (o *typ.Object, err error) {
				assert.Equal(t, srcPath, path)
				return &typ.Object{Name: srcPath}, nil
			})
			srcStore.EXPECT().String().DoAndReturn(func() string {
				return "src"
			}).AnyTimes()
			dstStore.EXPECT().Stat(gomock.Any()).DoAndReturn(func(path string, pairs ...*typ.Pair) (o *typ.Object, err error) {
				assert.Equal(t, dstPath, path)
				return tt.expectObject, tt.expectErr
			})
			dstStore.EXPECT().String().DoAndReturn(func() string {
				return "dst"
			}).AnyTimes()

			task.run()

			assert.NotNil(t, task.GetSourceObject())
			if tt.expectObject != nil {
				assert.NotNil(t, task.GetDestinationObject())
			} else {
				if tt.expectErr == typ.ErrObjectNotExist {
					assert.Nil(t, task.GetDestinationObject())
				} else {
					assert.Panics(t, func() {
						task.GetDestinationObject()
					})
				}
			}
		})
	}
}

func TestIsDestinationObjectExistTask_run(t *testing.T) {
	t.Run("destination object not exist", func(t *testing.T) {
		task := IsDestinationObjectExistTask{}
		task.SetDestinationObject(nil)

		task.run()

		assert.Equal(t, false, task.GetResult())
	})

	t.Run("destination object exists", func(t *testing.T) {
		task := IsDestinationObjectExistTask{}
		task.SetDestinationObject(&typ.Object{})

		task.run()

		assert.Equal(t, true, task.GetResult())
	})
}

func TestIsSizeEqualTask_run(t *testing.T) {
	t.Run("size equal", func(t *testing.T) {
		task := IsSizeEqualTask{}
		task.SetSourceObject(&typ.Object{Size: 111})
		task.SetDestinationObject(&typ.Object{Size: 111})

		task.run()

		assert.Equal(t, true, task.GetResult())
	})

	t.Run("size not equal", func(t *testing.T) {
		task := IsSizeEqualTask{}
		task.SetSourceObject(&typ.Object{Size: 222})
		task.SetDestinationObject(&typ.Object{Size: 111})

		task.run()

		assert.Equal(t, false, task.GetResult())
	})
}

func TestIsUpdateAtGreaterTask_run(t *testing.T) {
	t.Run("updated at greater", func(t *testing.T) {
		task := IsUpdateAtGreaterTask{}
		task.SetSourceObject(&typ.Object{UpdatedAt: time.Now().Add(time.Hour)})
		task.SetDestinationObject(&typ.Object{UpdatedAt: time.Now()})

		task.run()

		assert.Equal(t, true, task.GetResult())
	})

	t.Run("updated at not greater", func(t *testing.T) {
		task := IsUpdateAtGreaterTask{}
		task.SetSourceObject(&typ.Object{UpdatedAt: time.Now()})
		task.SetDestinationObject(&typ.Object{UpdatedAt: time.Now().Add(time.Hour)})

		task.run()

		assert.Equal(t, false, task.GetResult())
	})
}
