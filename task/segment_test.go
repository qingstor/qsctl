package task

import (
	"errors"
	"testing"

	typ "github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/mock"
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
