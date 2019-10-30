// Code generated by go generate; DO NOT EDIT.
package task

import (
	"errors"
	"testing"

	"github.com/Xuanwo/navvy"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool

func TestCopyFileTask_TriggerFault(t *testing.T) {
	m := &mockCopyFileTask{}
	task := &CopyFileTask{copyFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCopyFileTask_Run(t *testing.T) {
	task := &mockCopyFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestCopyLargeFileTask_TriggerFault(t *testing.T) {
	m := &mockCopyLargeFileTask{}
	task := &CopyLargeFileTask{copyLargeFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCopyLargeFileTask_Run(t *testing.T) {
	task := &mockCopyLargeFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestCopyPartialFileTask_TriggerFault(t *testing.T) {
	m := &mockCopyPartialFileTask{}
	task := &CopyPartialFileTask{copyPartialFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCopyPartialFileTask_Run(t *testing.T) {
	task := &mockCopyPartialFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestCopyPartialStreamTask_TriggerFault(t *testing.T) {
	m := &mockCopyPartialStreamTask{}
	task := &CopyPartialStreamTask{copyPartialStreamTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCopyPartialStreamTask_Run(t *testing.T) {
	task := &mockCopyPartialStreamTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestCopySmallFileTask_TriggerFault(t *testing.T) {
	m := &mockCopySmallFileTask{}
	task := &CopySmallFileTask{copySmallFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCopySmallFileTask_Run(t *testing.T) {
	task := &mockCopySmallFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestCopyStreamTask_TriggerFault(t *testing.T) {
	m := &mockCopyStreamTask{}
	task := &CopyStreamTask{copyStreamTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCopyStreamTask_Run(t *testing.T) {
	task := &mockCopyStreamTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestCreateStorageTask_TriggerFault(t *testing.T) {
	m := &mockCreateStorageTask{}
	task := &CreateStorageTask{createStorageTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockCreateStorageTask_Run(t *testing.T) {
	task := &mockCreateStorageTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestDeleteDirTask_TriggerFault(t *testing.T) {
	m := &mockDeleteDirTask{}
	task := &DeleteDirTask{deleteDirTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockDeleteDirTask_Run(t *testing.T) {
	task := &mockDeleteDirTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestDeleteFileTask_TriggerFault(t *testing.T) {
	m := &mockDeleteFileTask{}
	task := &DeleteFileTask{deleteFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockDeleteFileTask_Run(t *testing.T) {
	task := &mockDeleteFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestDeleteStorageTask_TriggerFault(t *testing.T) {
	m := &mockDeleteStorageTask{}
	task := &DeleteStorageTask{deleteStorageTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockDeleteStorageTask_Run(t *testing.T) {
	task := &mockDeleteStorageTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestDeleteStorageForceTask_TriggerFault(t *testing.T) {
	m := &mockDeleteStorageForceTask{}
	task := &DeleteStorageForceTask{deleteStorageForceTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockDeleteStorageForceTask_Run(t *testing.T) {
	task := &mockDeleteStorageForceTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestFileCopyTask_TriggerFault(t *testing.T) {
	m := &mockFileCopyTask{}
	task := &FileCopyTask{fileCopyTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockFileCopyTask_Run(t *testing.T) {
	task := &mockFileCopyTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestFileMD5SumTask_TriggerFault(t *testing.T) {
	m := &mockFileMD5SumTask{}
	task := &FileMD5SumTask{fileMD5SumTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockFileMD5SumTask_Run(t *testing.T) {
	task := &mockFileMD5SumTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestIterateFileTask_TriggerFault(t *testing.T) {
	m := &mockIterateFileTask{}
	task := &IterateFileTask{iterateFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockIterateFileTask_Run(t *testing.T) {
	task := &mockIterateFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestListFileTask_TriggerFault(t *testing.T) {
	m := &mockListFileTask{}
	task := &ListFileTask{listFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockListFileTask_Run(t *testing.T) {
	task := &mockListFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestListStorageTask_TriggerFault(t *testing.T) {
	m := &mockListStorageTask{}
	task := &ListStorageTask{listStorageTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockListStorageTask_Run(t *testing.T) {
	task := &mockListStorageTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestReachFileTask_TriggerFault(t *testing.T) {
	m := &mockReachFileTask{}
	task := &ReachFileTask{reachFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockReachFileTask_Run(t *testing.T) {
	task := &mockReachFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestSegmentAbortAllTask_TriggerFault(t *testing.T) {
	m := &mockSegmentAbortAllTask{}
	task := &SegmentAbortAllTask{segmentAbortAllTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockSegmentAbortAllTask_Run(t *testing.T) {
	task := &mockSegmentAbortAllTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestSegmentCompleteTask_TriggerFault(t *testing.T) {
	m := &mockSegmentCompleteTask{}
	task := &SegmentCompleteTask{segmentCompleteTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockSegmentCompleteTask_Run(t *testing.T) {
	task := &mockSegmentCompleteTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestSegmentFileCopyTask_TriggerFault(t *testing.T) {
	m := &mockSegmentFileCopyTask{}
	task := &SegmentFileCopyTask{segmentFileCopyTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockSegmentFileCopyTask_Run(t *testing.T) {
	task := &mockSegmentFileCopyTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestSegmentInitTask_TriggerFault(t *testing.T) {
	m := &mockSegmentInitTask{}
	task := &SegmentInitTask{segmentInitTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockSegmentInitTask_Run(t *testing.T) {
	task := &mockSegmentInitTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestSegmentStreamCopyTask_TriggerFault(t *testing.T) {
	m := &mockSegmentStreamCopyTask{}
	task := &SegmentStreamCopyTask{segmentStreamCopyTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockSegmentStreamCopyTask_Run(t *testing.T) {
	task := &mockSegmentStreamCopyTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestStatFileTask_TriggerFault(t *testing.T) {
	m := &mockStatFileTask{}
	task := &StatFileTask{statFileTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockStatFileTask_Run(t *testing.T) {
	task := &mockStatFileTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}

func TestStreamMD5SumTask_TriggerFault(t *testing.T) {
	m := &mockStreamMD5SumTask{}
	task := &StreamMD5SumTask{streamMD5SumTaskRequirement: m}
	err := errors.New("test error")
	task.TriggerFault(err)
	assert.True(t, task.GetFault().HasError())
}

func TestMockStreamMD5SumTask_Run(t *testing.T) {
	task := &mockStreamMD5SumTask{}
	assert.Panics(t, func() {
		task.Run()
	})
}
