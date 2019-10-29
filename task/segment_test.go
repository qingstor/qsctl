package task

import (
	"errors"
	"io"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/pkg/segment"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
	itypes "github.com/yunify/qsctl/v2/pkg/types"
	"github.com/yunify/qsctl/v2/utils"
)

func TestMultipartInitTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	x := &mockMultipartInitTask{}

	store := mock.NewMockStorager(ctrl)
	x.SetDestinationStorage(store)

	pool := navvy.NewPool(10)
	x.SetPool(pool)

	key := uuid.New().String()
	x.SetDestinationPath(key)

	offset := int64(0)
	x.SetCurrentOffset(&offset)
	x.SetTotalSize(1024)

	id := uuid.New().String()

	fn := func(task itypes.Todoist) navvy.Task {
		s := int64(1024)
		x.SetCurrentOffset(&s)

		t := &utils.EmptyTask{}
		t.SetID(id)
		t.SetPool(pool)
		return t
	}
	x.SetScheduler(itypes.NewScheduler(fn))

	store.EXPECT().InitSegment(gomock.Any()).Do(func(inputPath string) {
		assert.Equal(t, key, inputPath)
	})

	task := NewMultipartInitTask(x)
	task.Run()
}

func TestMultipartFileUploadTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := uuid.New().String()
	name := uuid.New().String()
	segmentID := uuid.New().String()
	size := int64(1234)

	srcStore := mock.NewMockStorager(ctrl)
	dstStore := mock.NewMockStorager(ctrl)

	x := &mockMultipartFileUploadTask{}
	x.SetDestinationPath(key)
	x.SetDestinationStorage(dstStore)
	x.SetSourcePath(name)
	x.SetSourceStorage(srcStore)
	x.SetSegmentID(segmentID)
	x.SetOffset(0)
	x.SetSize(size)
	x.SetID(uuid.New().String())

	sche := itypes.NewMockScheduler(nil)
	sche.New(nil)
	x.SetScheduler(sche)

	mockReader := mock.NewMockReadCloser(ctrl)
	mockReader.EXPECT().Close().Do(func() {})

	srcStore.EXPECT().Read(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, pairs ...*types.Pair) (r io.ReadCloser, err error) {
		assert.Equal(t, name, inputPath)
		assert.Equal(t, size, pairs[0].Value.(int64))
		assert.Equal(t, int64(0), pairs[1].Value.(int64))
		return mockReader, nil
	})
	dstStore.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, segmentID, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})

	task := NewMultipartFileUploadTask(x)
	task.Run()
}

func TestMultipartStreamUploadTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	x := &mockMultipartStreamUploadTask{}

	store := mock.NewMockStorager(ctrl)
	x.SetDestinationStorage(store)

	segmentID := uuid.New().String()
	x.SetSegmentID(segmentID)

	key := uuid.New().String()
	x.SetDestinationPath(key)

	buf, size, md5sum := utils.GenerateTestStream()

	x.SetSize(size)
	x.SetContent(buf)
	x.SetMD5Sum(md5sum)
	x.SetID(uuid.New().String())
	x.SetOffset(0)

	sche := itypes.NewMockScheduler(nil)
	sche.New(nil)
	x.SetScheduler(sche)

	store.EXPECT().WriteSegment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Do(func(inputPath string, inputOffset, inputSize int64, _ io.ReadCloser) {
		assert.Equal(t, segmentID, inputPath)
		assert.Equal(t, int64(0), inputOffset)
		assert.Equal(t, size, inputSize)
	})

	task := NewMultipartStreamUploadTask(x)
	task.Run()
}

func TestMultipartCompleteTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	x := &mockMultipartCompleteTask{}

	store := mock.NewMockStorager(ctrl)
	x.SetDestinationStorage(store)
	key := uuid.New().String()
	x.SetDestinationPath(key)
	segmentID := uuid.New().String()
	x.SetSegmentID(segmentID)

	store.EXPECT().CompleteSegment(gomock.Any()).Do(func(inputPath string, option ...*types.Pair) {
		assert.Equal(t, segmentID, inputPath)
	})

	task := NewMultipartCompleteTask(x)
	task.Run()
}

func TestAbortMultipartTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	bucketName := uuid.New().String()
	key, abortErr, listErr := "", errors.New("abort-multipart-abortErr"), errors.New("abort-multipart-listErr")
	pool := navvy.NewPool(10)

	cases := []struct {
		name       string
		listFault  bool
		listErr    error
		abortFault bool
		abortErr   error
	}{
		{
			name:       "ok",
			listFault:  false,
			listErr:    iterator.ErrDone,
			abortFault: false,
			abortErr:   nil,
		},
		{
			name:       "list error",
			listFault:  true,
			listErr:    listErr,
			abortFault: false,
			abortErr:   nil,
		},
		{
			name:       "abort error",
			listFault:  false,
			listErr:    iterator.ErrDone,
			abortFault: true,
			abortErr:   abortErr,
		},
	}

	for _, ca := range cases {
		x := &mockAbortMultipartTask{}
		x.SetBucketName(bucketName)
		x.SetPool(pool)

		store.EXPECT().ListSegments(gomock.Any()).DoAndReturn(func(inputPath string) iterator.SegmentIterator {
			assert.Equal(t, inputPath, key)

			count := 3
			return iterator.NewSegmentIterator(func(segments *[]*segment.Segment) error {
				*segments = make([]*segment.Segment, 1)
				(*segments)[0] = &segment.Segment{ID: ca.name}
				count--
				if count > 0 {
					return nil
				}
				return ca.listErr
			})
		})

		store.EXPECT().AbortSegment(gomock.Any()).DoAndReturn(func(inputPath string) error {
			return ca.abortErr
		}).AnyTimes()

		x.SetDestinationStorage(store)

		task := NewAbortMultipartTask(x)
		task.Run()

		assert.Equal(t, ca.listFault || ca.abortFault, x.ValidateFault(), ca.name)
		if ca.listFault {
			assert.Error(t, x.GetFault(), ca.name)
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.listErr), ca.name)
		}
		if ca.abortFault {
			assert.Error(t, x.GetFault(), ca.name)
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.abortErr), ca.name)
		}
	}
}
