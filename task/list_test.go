package task

import (
	"testing"

	"github.com/Xuanwo/storage/pkg/segment"
	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/pkg/fault"
	"github.com/qingstor/qsctl/v2/pkg/mock"
)

func TestListDirTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	testPath := uuid.New().String()

	task := ListDirTask{}
	task.SetFault(fault.New())
	task.SetStorage(store)
	task.SetPath(testPath)
	task.SetDirFunc(func(*typ.Object) {})
	task.SetFileFunc(func(*typ.Object) {})

	store.EXPECT().ListDir(gomock.Any(), gomock.Any()).Do(func(path string, opts ...*typ.Pair) error {
		assert.Equal(t, testPath, path)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestListSegmentTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	segmenter := mock.NewMockSegmenter(ctrl)
	testPath := uuid.New().String()

	task := ListSegmentTask{}
	task.SetFault(fault.New())
	task.SetSegmenter(segmenter)
	task.SetPath(testPath)
	task.SetSegmentFunc(func(segment *segment.Segment) {})

	segmenter.EXPECT().ListSegments(gomock.Any(), gomock.Any()).Do(func(path string, opts ...*typ.Pair) error {
		assert.Equal(t, testPath, path)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestListStorageTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := mock.NewMockServicer(ctrl)
	zone := uuid.New().String()

	task := ListStorageTask{}
	task.SetFault(fault.New())
	task.SetService(srv)
	task.SetZone(zone)

	srv.EXPECT().List(gomock.Any()).Do(func(pairs ...*typ.Pair) error {
		assert.Equal(t, zone, pairs[0].Value.(string))
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}
