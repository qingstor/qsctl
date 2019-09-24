// Code generated by go generate; DO NOT EDIT.
package types

import (
	"bytes"
	"io"
	"sync"

	"github.com/Xuanwo/navvy"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
)

type BucketNameGetter interface {
	GetBucketName() string
}

type BucketNameSetter interface {
	SetBucketName(string)
}

type BucketName struct {
	valid bool
	v     string
}

func (o *BucketName) GetBucketName() string {
	if !o.valid {
		panic("BucketName value is not valid")
	}
	return o.v
}

func (o *BucketName) SetBucketName(v string) {
	o.v = v
	o.valid = true
}

type BytesPoolGetter interface {
	GetBytesPool() *sync.Pool
}

type BytesPoolSetter interface {
	SetBytesPool(*sync.Pool)
}

type BytesPool struct {
	valid bool
	v     *sync.Pool
}

func (o *BytesPool) GetBytesPool() *sync.Pool {
	if !o.valid {
		panic("BytesPool value is not valid")
	}
	return o.v
}

func (o *BytesPool) SetBytesPool(v *sync.Pool) {
	o.v = v
	o.valid = true
}

type ContentGetter interface {
	GetContent() *bytes.Buffer
}

type ContentSetter interface {
	SetContent(*bytes.Buffer)
}

type Content struct {
	valid bool
	v     *bytes.Buffer
}

func (o *Content) GetContent() *bytes.Buffer {
	if !o.valid {
		panic("Content value is not valid")
	}
	return o.v
}

func (o *Content) SetContent(v *bytes.Buffer) {
	o.v = v
	o.valid = true
}

type CurrentOffsetGetter interface {
	GetCurrentOffset() *int64
}

type CurrentOffsetSetter interface {
	SetCurrentOffset(*int64)
}

type CurrentOffset struct {
	valid bool
	v     *int64
}

func (o *CurrentOffset) GetCurrentOffset() *int64 {
	if !o.valid {
		panic("CurrentOffset value is not valid")
	}
	return o.v
}

func (o *CurrentOffset) SetCurrentOffset(v *int64) {
	o.v = v
	o.valid = true
}

type CurrentPartNumberGetter interface {
	GetCurrentPartNumber() *int32
}

type CurrentPartNumberSetter interface {
	SetCurrentPartNumber(*int32)
}

type CurrentPartNumber struct {
	valid bool
	v     *int32
}

func (o *CurrentPartNumber) GetCurrentPartNumber() *int32 {
	if !o.valid {
		panic("CurrentPartNumber value is not valid")
	}
	return o.v
}

func (o *CurrentPartNumber) SetCurrentPartNumber(v *int32) {
	o.v = v
	o.valid = true
}

type EnableBenchmarkGetter interface {
	GetEnableBenchmark() bool
}

type EnableBenchmarkSetter interface {
	SetEnableBenchmark(bool)
}

type EnableBenchmark struct {
	valid bool
	v     bool
}

func (o *EnableBenchmark) GetEnableBenchmark() bool {
	if !o.valid {
		panic("EnableBenchmark value is not valid")
	}
	return o.v
}

func (o *EnableBenchmark) SetEnableBenchmark(v bool) {
	o.v = v
	o.valid = true
}

type ExpectSizeGetter interface {
	GetExpectSize() int64
}

type ExpectSizeSetter interface {
	SetExpectSize(int64)
}

type ExpectSize struct {
	valid bool
	v     int64
}

func (o *ExpectSize) GetExpectSize() int64 {
	if !o.valid {
		panic("ExpectSize value is not valid")
	}
	return o.v
}

func (o *ExpectSize) SetExpectSize(v int64) {
	o.v = v
	o.valid = true
}

type FlowTypeGetter interface {
	GetFlowType() constants.FlowType
}

type FlowTypeSetter interface {
	SetFlowType(constants.FlowType)
}

type FlowType struct {
	valid bool
	v     constants.FlowType
}

func (o *FlowType) GetFlowType() constants.FlowType {
	if !o.valid {
		panic("FlowType value is not valid")
	}
	return o.v
}

func (o *FlowType) SetFlowType(v constants.FlowType) {
	o.v = v
	o.valid = true
}

type FormatGetter interface {
	GetFormat() string
}

type FormatSetter interface {
	SetFormat(string)
}

type Format struct {
	valid bool
	v     string
}

func (o *Format) GetFormat() string {
	if !o.valid {
		panic("Format value is not valid")
	}
	return o.v
}

func (o *Format) SetFormat(v string) {
	o.v = v
	o.valid = true
}

type KeyGetter interface {
	GetKey() string
}

type KeySetter interface {
	SetKey(string)
}

type Key struct {
	valid bool
	v     string
}

func (o *Key) GetKey() string {
	if !o.valid {
		panic("Key value is not valid")
	}
	return o.v
}

func (o *Key) SetKey(v string) {
	o.v = v
	o.valid = true
}

type KeyTypeGetter interface {
	GetKeyType() constants.KeyType
}

type KeyTypeSetter interface {
	SetKeyType(constants.KeyType)
}

type KeyType struct {
	valid bool
	v     constants.KeyType
}

func (o *KeyType) GetKeyType() constants.KeyType {
	if !o.valid {
		panic("KeyType value is not valid")
	}
	return o.v
}

func (o *KeyType) SetKeyType(v constants.KeyType) {
	o.v = v
	o.valid = true
}

type MD5SumGetter interface {
	GetMD5Sum() []byte
}

type MD5SumSetter interface {
	SetMD5Sum([]byte)
}

type MD5Sum struct {
	valid bool
	v     []byte
}

func (o *MD5Sum) GetMD5Sum() []byte {
	if !o.valid {
		panic("MD5Sum value is not valid")
	}
	return o.v
}

func (o *MD5Sum) SetMD5Sum(v []byte) {
	o.v = v
	o.valid = true
}

type ObjectMetaGetter interface {
	GetObjectMeta() *storage.ObjectMeta
}

type ObjectMetaSetter interface {
	SetObjectMeta(*storage.ObjectMeta)
}

type ObjectMeta struct {
	valid bool
	v     *storage.ObjectMeta
}

func (o *ObjectMeta) GetObjectMeta() *storage.ObjectMeta {
	if !o.valid {
		panic("ObjectMeta value is not valid")
	}
	return o.v
}

func (o *ObjectMeta) SetObjectMeta(v *storage.ObjectMeta) {
	o.v = v
	o.valid = true
}

type OffsetGetter interface {
	GetOffset() int64
}

type OffsetSetter interface {
	SetOffset(int64)
}

type Offset struct {
	valid bool
	v     int64
}

func (o *Offset) GetOffset() int64 {
	if !o.valid {
		panic("Offset value is not valid")
	}
	return o.v
}

func (o *Offset) SetOffset(v int64) {
	o.v = v
	o.valid = true
}

type PartNumberGetter interface {
	GetPartNumber() int
}

type PartNumberSetter interface {
	SetPartNumber(int)
}

type PartNumber struct {
	valid bool
	v     int
}

func (o *PartNumber) GetPartNumber() int {
	if !o.valid {
		panic("PartNumber value is not valid")
	}
	return o.v
}

func (o *PartNumber) SetPartNumber(v int) {
	o.v = v
	o.valid = true
}

type PartSizeGetter interface {
	GetPartSize() int64
}

type PartSizeSetter interface {
	SetPartSize(int64)
}

type PartSize struct {
	valid bool
	v     int64
}

func (o *PartSize) GetPartSize() int64 {
	if !o.valid {
		panic("PartSize value is not valid")
	}
	return o.v
}

func (o *PartSize) SetPartSize(v int64) {
	o.v = v
	o.valid = true
}

type PathGetter interface {
	GetPath() string
}

type PathSetter interface {
	SetPath(string)
}

type Path struct {
	valid bool
	v     string
}

func (o *Path) GetPath() string {
	if !o.valid {
		panic("Path value is not valid")
	}
	return o.v
}

func (o *Path) SetPath(v string) {
	o.v = v
	o.valid = true
}

type PathTypeGetter interface {
	GetPathType() constants.PathType
}

type PathTypeSetter interface {
	SetPathType(constants.PathType)
}

type PathType struct {
	valid bool
	v     constants.PathType
}

func (o *PathType) GetPathType() constants.PathType {
	if !o.valid {
		panic("PathType value is not valid")
	}
	return o.v
}

func (o *PathType) SetPathType(v constants.PathType) {
	o.v = v
	o.valid = true
}

type PoolGetter interface {
	GetPool() *navvy.Pool
}

type PoolSetter interface {
	SetPool(*navvy.Pool)
}

type Pool struct {
	valid bool
	v     *navvy.Pool
}

func (o *Pool) GetPool() *navvy.Pool {
	if !o.valid {
		panic("Pool value is not valid")
	}
	return o.v
}

func (o *Pool) SetPool(v *navvy.Pool) {
	o.v = v
	o.valid = true
}

type SizeGetter interface {
	GetSize() int64
}

type SizeSetter interface {
	SetSize(int64)
}

type Size struct {
	valid bool
	v     int64
}

func (o *Size) GetSize() int64 {
	if !o.valid {
		panic("Size value is not valid")
	}
	return o.v
}

func (o *Size) SetSize(v int64) {
	o.v = v
	o.valid = true
}

type StorageGetter interface {
	GetStorage() storage.ObjectStorage
}

type StorageSetter interface {
	SetStorage(storage.ObjectStorage)
}

type Storage struct {
	valid bool
	v     storage.ObjectStorage
}

func (o *Storage) GetStorage() storage.ObjectStorage {
	if !o.valid {
		panic("Storage value is not valid")
	}
	return o.v
}

func (o *Storage) SetStorage(v storage.ObjectStorage) {
	o.v = v
	o.valid = true
}

type StreamGetter interface {
	GetStream() io.Reader
}

type StreamSetter interface {
	SetStream(io.Reader)
}

type Stream struct {
	valid bool
	v     io.Reader
}

func (o *Stream) GetStream() io.Reader {
	if !o.valid {
		panic("Stream value is not valid")
	}
	return o.v
}

func (o *Stream) SetStream(v io.Reader) {
	o.v = v
	o.valid = true
}

type TaskConstructorGetter interface {
	GetTaskConstructor() TodoFunc
}

type TaskConstructorSetter interface {
	SetTaskConstructor(TodoFunc)
}

type TaskConstructor struct {
	valid bool
	v     TodoFunc
}

func (o *TaskConstructor) GetTaskConstructor() TodoFunc {
	if !o.valid {
		panic("TaskConstructor value is not valid")
	}
	return o.v
}

func (o *TaskConstructor) SetTaskConstructor(v TodoFunc) {
	o.v = v
	o.valid = true
}

type TotalSizeGetter interface {
	GetTotalSize() int64
}

type TotalSizeSetter interface {
	SetTotalSize(int64)
}

type TotalSize struct {
	valid bool
	v     int64
}

func (o *TotalSize) GetTotalSize() int64 {
	if !o.valid {
		panic("TotalSize value is not valid")
	}
	return o.v
}

func (o *TotalSize) SetTotalSize(v int64) {
	o.v = v
	o.valid = true
}

type UploadIDGetter interface {
	GetUploadID() string
}

type UploadIDSetter interface {
	SetUploadID(string)
}

type UploadID struct {
	valid bool
	v     string
}

func (o *UploadID) GetUploadID() string {
	if !o.valid {
		panic("UploadID value is not valid")
	}
	return o.v
}

func (o *UploadID) SetUploadID(v string) {
	o.v = v
	o.valid = true
}

type WaitGroupGetter interface {
	GetWaitGroup() *sync.WaitGroup
}

type WaitGroupSetter interface {
	SetWaitGroup(*sync.WaitGroup)
}

type WaitGroup struct {
	valid bool
	v     *sync.WaitGroup
}

func (o *WaitGroup) GetWaitGroup() *sync.WaitGroup {
	if !o.valid {
		panic("WaitGroup value is not valid")
	}
	return o.v
}

func (o *WaitGroup) SetWaitGroup(v *sync.WaitGroup) {
	o.v = v
	o.valid = true
}

type ZoneGetter interface {
	GetZone() string
}

type ZoneSetter interface {
	SetZone(string)
}

type Zone struct {
	valid bool
	v     string
}

func (o *Zone) GetZone() string {
	if !o.valid {
		panic("Zone value is not valid")
	}
	return o.v
}

func (o *Zone) SetZone(v string) {
	o.v = v
	o.valid = true
}
