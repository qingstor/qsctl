package common

import (
	"bufio"
	"crypto/md5"
	"io"
	"os"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"

	"github.com/yunify/qsctl/v2/task/types"
	"github.com/yunify/qsctl/v2/task/utils"
)

// FileMD5SumTaskRequirement is the requirement for execute FileMD5SumTask.
type FileMD5SumTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.MD5SumSetter
	types.PathGetter
	types.OffsetGetter
	types.SizeGetter
	types.PoolGetter
}

// FileMD5SumTask will execute SeekableMD5Sum Task.
type FileMD5SumTask struct {
	FileMD5SumTaskRequirement
}

// NewFileMD5SumTask will create a new Task.
func NewFileMD5SumTask(task types.Todoist) navvy.Task {
	o, ok := task.(FileMD5SumTaskRequirement)
	if !ok {
		panic("task is not fill FileMD5SumTaskRequirement")
	}

	return &FileMD5SumTask{o}
}

// Run implement navvy.Task.
func (t *FileMD5SumTask) Run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "FileMD5SumTask", t.GetPath(), t.GetOffset())

	f, err := os.Open(t.GetPath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(io.NewSectionReader(f, t.GetOffset(), t.GetSize()))
	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		log.Errorf("Task <%s> failed for [%v]", "FileMD5SumTask", err)
	}

	t.SetMD5Sum(h.Sum(nil)[:])

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "FileMD5SumTask", t.GetPath(), t.GetOffset())
	utils.SubmitNextTask(t.FileMD5SumTaskRequirement)
}

// StreamMD5SumTaskRequirement is the requirement for execute StreamMD5SumTask.
type StreamMD5SumTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.MD5SumSetter
	types.PathGetter
	types.ContentGetter
	types.PoolGetter
}

// StreamMD5SumTask will execute SeekableMD5Sum Task.
type StreamMD5SumTask struct {
	StreamMD5SumTaskRequirement
}

// NewStreamMD5SumTask will create a new Task.
func NewStreamMD5SumTask(task types.Todoist) navvy.Task {
	o, ok := task.(StreamMD5SumTaskRequirement)
	if !ok {
		panic("task is not fill StreamMD5SumTaskRequirement")
	}

	return &StreamMD5SumTask{o}
}

// Run implement navvy.Task.
func (t *StreamMD5SumTask) Run() {
	log.Debugf("Task <%s> for Stream <%s> started.", "StreamMD5SumTask", t.GetPath())

	md5Sum := md5.Sum(t.GetContent().Bytes())
	t.SetMD5Sum(md5Sum[:])

	log.Debugf("Task <%s> for Stream <%s> finished.", "StreamMD5SumTask", t.GetPath())
	utils.SubmitNextTask(t.StreamMD5SumTaskRequirement)
}
