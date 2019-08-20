package common

import (
	"bufio"
	"crypto/md5"
	"io"
	"os"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qsctl/v2/task/utils"

	"github.com/yunify/qsctl/v2/task/types"
)

// SeekableMD5SumTaskRequirement is the requirement for execute SeekableMD5SumTask.
type SeekableMD5SumTaskRequirement interface {
	navvy.Task
	types.Todoist

	types.MD5SumSetter
	types.FilePathGetter
	types.OffsetGetter
	types.SizeGetter
	types.PoolGetter
}

// SeekableMD5SumTask will execute SeekableMD5Sum Task.
type SeekableMD5SumTask struct {
	SeekableMD5SumTaskRequirement
}

// NewSeekableMD5SumTask will create a new Task.
func NewSeekableMD5SumTask(task types.Todoist) navvy.Task {
	o, ok := task.(SeekableMD5SumTaskRequirement)
	if !ok {
		panic("task is not fill SeekableMD5SumTaskRequirement")
	}

	return &SeekableMD5SumTask{o}
}

// Run implement navvy.Task.
func (t *SeekableMD5SumTask) Run() {
	log.Debugf("Task <%s> for File <%s> at Offset <%d> started.", "SeekableMD5SumTask", t.GetFilePath(), t.GetOffset())

	f, err := os.Open(t.GetFilePath())
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(io.NewSectionReader(f, t.GetOffset(), t.GetSize()))
	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		log.Errorf("Task <%s> failed for [%v]", "SeekableMD5SumTask", err)
	}

	t.SetMD5Sum(h.Sum(nil)[:])

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "SeekableMD5SumTask", t.GetFilePath(), t.GetOffset())
	go utils.SubmitNextTask(t.SeekableMD5SumTaskRequirement)
}
