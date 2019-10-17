// Code generated by go generate; DO NOT EDIT.
package fault

import (
	"fmt"

	"github.com/yunify/qsctl/v2/pkg/types"
)
type LocalFileNotExist struct {
	types.Fault
	types.Path
}

func (f *LocalFileNotExist) Error() string {
	return fmt.Sprintf(`Local file [%s] is not exist: {%v}`, f.GetPath(), f.GetFault())
}

func (f *LocalFileNotExist) Unwrap() error {
	return f.GetFault()
}

func NewLocalFileNotExist(err error,path string) error {
	f := &LocalFileNotExist{}
	f.SetFault(err)
	f.SetPath(path)

	return f
}
type LocalFileTooLarge struct {
	types.Fault
	types.Size
}

func (f *LocalFileTooLarge) Error() string {
	return fmt.Sprintf(`Local file size [%d] is too large`, f.GetSize())
}

func (f *LocalFileTooLarge) Unwrap() error {
	return f.GetFault()
}

func NewLocalFileTooLarge(err error,size int64) error {
	f := &LocalFileTooLarge{}
	f.SetFault(err)
	f.SetSize(size)

	return f
}
type ReadableSizeFormatInvalid struct {
	types.Fault
	types.ByteSize
}

func (f *ReadableSizeFormatInvalid) Error() string {
	return fmt.Sprintf(`readable size format invalid [%s]`, f.GetByteSize())
}

func (f *ReadableSizeFormatInvalid) Unwrap() error {
	return f.GetFault()
}

func NewReadableSizeFormatInvalid(err error,byteSize string) error {
	f := &ReadableSizeFormatInvalid{}
	f.SetFault(err)
	f.SetByteSize(byteSize)

	return f
}
type StorageBucketInitFailed struct {
	types.Fault
	types.BucketName
	types.Zone
}

func (f *StorageBucketInitFailed) Error() string {
	return fmt.Sprintf(`Storage bucket [%s] in zone [%s] initiate failed: {%v}`, f.GetBucketName(), f.GetZone(), f.GetFault())
}

func (f *StorageBucketInitFailed) Unwrap() error {
	return f.GetFault()
}

func NewStorageBucketInitFailed(err error,bucketName string,zone string) error {
	f := &StorageBucketInitFailed{}
	f.SetFault(err)
	f.SetBucketName(bucketName)
	f.SetZone(zone)

	return f
}
type StorageObjectNoPermission struct {
	types.Fault
	types.Key
}

func (f *StorageObjectNoPermission) Error() string {
	return fmt.Sprintf(`Storage Object [%s] do not have enough permission: {%v}`, f.GetKey(), f.GetFault())
}

func (f *StorageObjectNoPermission) Unwrap() error {
	return f.GetFault()
}

func NewStorageObjectNoPermission(err error,key string) error {
	f := &StorageObjectNoPermission{}
	f.SetFault(err)
	f.SetKey(key)

	return f
}
type StorageObjectNotFound struct {
	types.Fault
	types.Key
}

func (f *StorageObjectNotFound) Error() string {
	return fmt.Sprintf(`Storage Object [%s] is not found: {%v}`, f.GetKey(), f.GetFault())
}

func (f *StorageObjectNotFound) Unwrap() error {
	return f.GetFault()
}

func NewStorageObjectNotFound(err error,key string) error {
	f := &StorageObjectNotFound{}
	f.SetFault(err)
	f.SetKey(key)

	return f
}
type StorageServiceInitFailed struct {
	types.Fault
}

func (f *StorageServiceInitFailed) Error() string {
	return fmt.Sprintf(`Storage service initiate failed: {%v}`, f.GetFault())
}

func (f *StorageServiceInitFailed) Unwrap() error {
	return f.GetFault()
}

func NewStorageServiceInitFailed(err error) error {
	f := &StorageServiceInitFailed{}
	f.SetFault(err)

	return f
}
type Unhandled struct {
	types.Fault
}

func (f *Unhandled) Error() string {
	return fmt.Sprintf(`Operation failed via unhandled error: {%v}`, f.GetFault())
}

func (f *Unhandled) Unwrap() error {
	return f.GetFault()
}

func NewUnhandled(err error) error {
	f := &Unhandled{}
	f.SetFault(err)

	return f
}
type UserInputByteSizeInvalid struct {
	types.Fault
	types.ByteSize
}

func (f *UserInputByteSizeInvalid) Error() string {
	return fmt.Sprintf(`User input byte size [%s] is invalid: {%v}`, f.GetByteSize(), f.GetFault())
}

func (f *UserInputByteSizeInvalid) Unwrap() error {
	return f.GetFault()
}

func NewUserInputByteSizeInvalid(err error,byteSize string) error {
	f := &UserInputByteSizeInvalid{}
	f.SetFault(err)
	f.SetByteSize(byteSize)

	return f
}
type UserInputKeyInvalid struct {
	types.Fault
	types.Key
}

func (f *UserInputKeyInvalid) Error() string {
	return fmt.Sprintf(`User input key [%s] is invalid`, f.GetKey())
}

func (f *UserInputKeyInvalid) Unwrap() error {
	return f.GetFault()
}

func NewUserInputKeyInvalid(err error,key string) error {
	f := &UserInputKeyInvalid{}
	f.SetFault(err)
	f.SetKey(key)

	return f
}
