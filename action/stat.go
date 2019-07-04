package action

import (
	"github.com/c2h5oh/datasize"
	"github.com/jedib0t/go-pretty/text"
	"github.com/yunify/qsctl/helper"
	"strings"
)

// Stat will handle all stat actions.
func Stat(remote string) (err error) {
	objectKey, err := ParseQsPath(remote)
	if err != nil {
		panic(err)
	}

	return StatRemoteObject(objectKey)
}

// StatRemoteObject will stat a remote object.
func StatRemoteObject(objectKey string) (err error) {
	om, err := helper.HeadObject(objectKey)
	if err != nil {
		panic(err)
	}

	content := []string{
		text.AlignLeft.Apply("Key: ", 14) + text.AlignLeft.Apply(om.Key, 1),
		text.AlignLeft.Apply("Size: ", 14) + text.AlignLeft.Apply(datasize.ByteSize(om.ContentLength).String(), 1),
		text.AlignLeft.Apply("Type: ", 14) + text.AlignLeft.Apply(om.ContentType, 1),
		text.AlignLeft.Apply("Modify: ", 14) + text.AlignLeft.Apply(om.LastModified.String(), 1),
		text.AlignLeft.Apply("StorageClass: ", 14) + text.AlignLeft.Apply(om.StorageClass, 1),
	}

	if om.ETag != "" {
		content = append(content, text.AlignLeft.Apply("MD5: ", 14)+text.AlignLeft.Apply(om.ETag, 1))
	}

	println(strings.Join(content, "\n"))
	return
}
