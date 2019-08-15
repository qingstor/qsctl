package task

import (
	"bufio"
	"crypto/md5"
	"io"
	"os"

	"github.com/Xuanwo/navvy"
	log "github.com/sirupsen/logrus"
)

// SeekableMD5SumTaskRequirement is the requirement for execute SeekableMD5SumTask.
type SeekableMD5SumTaskRequirement interface {
	navvy.Task
	Todoist

	MD5SumSetter
	FilePathGetter
	OffsetGetter
	ContentLengthGetter
}

// SeekableMD5SumTask will execute SeekableMD5Sum Task.
type SeekableMD5SumTask struct {
	SeekableMD5SumTaskRequirement
}

// NewSeekableMD5SumTask will create a new Task.
func NewSeekableMD5SumTask(task Todoist) navvy.Task {
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

	r := bufio.NewReader(io.NewSectionReader(f, t.GetOffset(), t.GetContentLength()))
	h := md5.New()
	_, err = io.Copy(h, r)
	if err != nil {
		log.Errorf("Task <%s> failed for [%v]", "SeekableMD5SumTask", err)
	}

	t.SetMD5Sum(h.Sum(nil)[:])

	log.Debugf("Task <%s> for File <%s> at Offset <%d> finished.", "SeekableMD5SumTask", t.GetFilePath(), t.GetOffset())
	go SubmitNextTask(t.SeekableMD5SumTaskRequirement)
}

// WaitTaskRequirement is the requirement for execute WaitTask.
type WaitTaskRequirement interface {
	navvy.Task
	Todoist

	WaitGroupGetter
}

// WaitTask will execute Wait Task.
type WaitTask struct {
	WaitTaskRequirement
}

// NewWaitTask will create a new Task.
func NewWaitTask(task Todoist) navvy.Task {
	o, ok := task.(WaitTaskRequirement)
	if !ok {
		panic("task is not fill NewWaitTask")
	}

	return &WaitTask{o}
}

// Run implement navvy.Task.
func (t *WaitTask) Run() {
	t.GetWaitGroup().Wait()

	go SubmitNextTask(t.WaitTaskRequirement)
}
