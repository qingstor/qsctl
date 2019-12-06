package task

import (
	"bytes"
	"errors"
	"io"
	"testing"

	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/qingstor/qsctl/v2/pkg/fault"
	"github.com/qingstor/qsctl/v2/pkg/mock"
)

func TestSegmentInitTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("normal", func(t *testing.T) {
		store := mock.NewMockSegmenter(ctrl)
		path := uuid.New().String()
		segmentID := uuid.New().String()

		task := SegmentInitTask{}
		task.SetSegmenter(store)
		task.SetPath(path)
		task.SetPartSize(1000)

		store.EXPECT().InitSegment(gomock.Any(), gomock.Any()).DoAndReturn(
			func(inputPath string, pairs ...*typ.Pair) (string, error) {
				assert.Equal(t, inputPath, path)
				assert.Equal(t, pairs[0].Value.(int64), int64(1000))
				return segmentID, nil
			},
		)

		task.run()

		assert.Equal(t, segmentID, task.GetSegmentID())
	})

	t.Run("init segment returned error", func(t *testing.T) {
		store := mock.NewMockSegmenter(ctrl)
		path := uuid.New().String()

		task := SegmentInitTask{}
		task.SetFault(fault.New())
		task.SetSegmenter(store)
		task.SetPath(path)
		task.SetPartSize(1000)

		store.EXPECT().InitSegment(gomock.Any(), gomock.Any()).DoAndReturn(
			func(inputPath string, pairs ...*typ.Pair) (string, error) {
				assert.Equal(t, inputPath, path)
				assert.Equal(t, pairs[0].Value.(int64), int64(1000))
				return "", errors.New("test")
			},
		)

		task.run()

		assert.False(t, task.ValidateSegmentID())
		assert.True(t, task.GetFault().HasError())
	})
}

func TestSegmentFileCopyTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srcStore := mock.NewMockStorager(ctrl)
	srcPath := uuid.New().String()
	srcReader := mock.NewMockReadCloser(ctrl)
	srcSize := int64(1024)
	dstSegmentID := uuid.New().String()
	dstSegmenter := mock.NewMockSegmenter(ctrl)
	dstPath := uuid.New().String()

	task := SegmentFileCopyTask{}
	task.SetFault(fault.New())
	task.SetSourcePath(srcPath)
	task.SetSourceStorage(srcStore)
	task.SetDestinationPath(dstPath)
	task.SetDestinationSegmenter(dstSegmenter)
	task.SetSize(srcSize)
	task.SetOffset(0)
	task.SetSegmentID(dstSegmentID)

	srcReader.EXPECT().Close()
	srcStore.EXPECT().Read(gomock.Any(), gomock.Any()).DoAndReturn(func(path string, pairs ...*typ.Pair) (r io.ReadCloser, err error) {
		assert.Equal(t, srcPath, path)
		return srcReader, nil
	})
	dstSegmenter.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(path string, offset, size int64, r io.Reader) (err error) {
		assert.Equal(t, dstSegmentID, path)
		assert.Equal(t, int64(0), offset)
		assert.Equal(t, srcSize, size)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestSegmentStreamCopyTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srcSize := int64(1024)
	dstSegmentID := uuid.New().String()
	segmenter := mock.NewMockSegmenter(ctrl)

	task := SegmentStreamCopyTask{}
	task.SetFault(fault.New())
	task.SetSegmenter(segmenter)
	task.SetSize(srcSize)
	task.SetOffset(0)
	task.SetSegmentID(dstSegmentID)
	task.SetContent(&bytes.Buffer{})

	segmenter.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(path string, offset, size int64, r io.Reader) (err error) {
		assert.Equal(t, dstSegmentID, path)
		assert.Equal(t, int64(0), offset)
		assert.Equal(t, srcSize, size)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}

func TestSegmentCompleteTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dstSegmentID := uuid.New().String()
	segmenter := mock.NewMockSegmenter(ctrl)

	task := SegmentCompleteTask{}
	task.SetFault(fault.New())
	task.SetSegmenter(segmenter)
	task.SetSegmentID(dstSegmentID)

	segmenter.EXPECT().CompleteSegment(gomock.Any()).DoAndReturn(func(path string) (err error) {
		assert.Equal(t, dstSegmentID, path)
		return nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}
