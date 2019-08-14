package task

// CopyLargeFileTask will execute CopyPartialFile Task
//
// CopyObject
// CopySingleFile
// CopyLargeFile
//   -> CopyPartialFile
// CopyNotSeekableFile
type CopyLargeFileTask struct {
	Todo

	FilePath
	ObjectKey
	UploadID
	TotalParts
	WaitGroup
}

// Run implement navvy.Task
func (t *CopyLargeFileTask) Run() {
	SubmitNextTask(t)
}

// CopyPartialFileTask will execute CopyPartialFile Task
type CopyPartialFileTask struct {
	Todo

	MD5Sum
	ContentLength
	PartNumber
	UploadID
	FilePath
	ObjectKey
	Offset
	WaitGroup
}

// NewCopyPartialFileTask will create a new Task.
func NewCopyPartialFileTask(
	objectKey, filePath, uploadID string,
	partNumber int,
	offset, contentLength int64,
) *CopyPartialFileTask {
	t := &CopyPartialFileTask{}
	t.SetPartNumber(partNumber)
	t.SetOffset(offset)
	t.SetContentLength(contentLength)
	t.SetUploadID(uploadID)
	t.SetFilePath(filePath)
	t.SetObjectKey(objectKey)

	t.AddTODOs(
		NewSeekableMD5SumTask,
		NewMultipartObjectUploadTask,
	)
	return t
}

// Run implement navvy.Task
func (t *CopyPartialFileTask) Run() {
	SubmitNextTask(t)
}
