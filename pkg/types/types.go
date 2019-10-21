// Code generated by go generate; DO NOT EDIT.
package types

import (
	"bytes"
	"io"
	"sync"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"

	"github.com/yunify/qsctl/v2/constants"
)

type BucketList struct {
	valid bool
	v     []string
}

type BucketListGetter interface {
	GetBucketList() []string
}

func (o *BucketList) GetBucketList() []string {
	if !o.valid {
		panic("BucketList value is not valid")
	}
	return o.v
}

type BucketListSetter interface {
	SetBucketList([]string)
}

func (o *BucketList) SetBucketList(v []string) {
	o.v = v
	o.valid = true
}

type BucketListValidator interface {
	ValidateBucketList() bool
}

func (o *BucketList) ValidateBucketList() bool {
	return o.valid
}

type BucketName struct {
	valid bool
	v     string
}

type BucketNameGetter interface {
	GetBucketName() string
}

func (o *BucketName) GetBucketName() string {
	if !o.valid {
		panic("BucketName value is not valid")
	}
	return o.v
}

type BucketNameSetter interface {
	SetBucketName(string)
}

func (o *BucketName) SetBucketName(v string) {
	o.v = v
	o.valid = true
}

type BucketNameValidator interface {
	ValidateBucketName() bool
}

func (o *BucketName) ValidateBucketName() bool {
	return o.valid
}

type ByteSize struct {
	valid bool
	v     string
}

type ByteSizeGetter interface {
	GetByteSize() string
}

func (o *ByteSize) GetByteSize() string {
	if !o.valid {
		panic("ByteSize value is not valid")
	}
	return o.v
}

type ByteSizeSetter interface {
	SetByteSize(string)
}

func (o *ByteSize) SetByteSize(v string) {
	o.v = v
	o.valid = true
}

type ByteSizeValidator interface {
	ValidateByteSize() bool
}

func (o *ByteSize) ValidateByteSize() bool {
	return o.valid
}

type BytesPool struct {
	valid bool
	v     *sync.Pool
}

type BytesPoolGetter interface {
	GetBytesPool() *sync.Pool
}

func (o *BytesPool) GetBytesPool() *sync.Pool {
	if !o.valid {
		panic("BytesPool value is not valid")
	}
	return o.v
}

type BytesPoolSetter interface {
	SetBytesPool(*sync.Pool)
}

func (o *BytesPool) SetBytesPool(v *sync.Pool) {
	o.v = v
	o.valid = true
}

type BytesPoolValidator interface {
	ValidateBytesPool() bool
}

func (o *BytesPool) ValidateBytesPool() bool {
	return o.valid
}

type Content struct {
	valid bool
	v     *bytes.Buffer
}

type ContentGetter interface {
	GetContent() *bytes.Buffer
}

func (o *Content) GetContent() *bytes.Buffer {
	if !o.valid {
		panic("Content value is not valid")
	}
	return o.v
}

type ContentSetter interface {
	SetContent(*bytes.Buffer)
}

func (o *Content) SetContent(v *bytes.Buffer) {
	o.v = v
	o.valid = true
}

type ContentValidator interface {
	ValidateContent() bool
}

func (o *Content) ValidateContent() bool {
	return o.valid
}

type CurrentOffset struct {
	valid bool
	v     *int64
}

type CurrentOffsetGetter interface {
	GetCurrentOffset() *int64
}

func (o *CurrentOffset) GetCurrentOffset() *int64 {
	if !o.valid {
		panic("CurrentOffset value is not valid")
	}
	return o.v
}

type CurrentOffsetSetter interface {
	SetCurrentOffset(*int64)
}

func (o *CurrentOffset) SetCurrentOffset(v *int64) {
	o.v = v
	o.valid = true
}

type CurrentOffsetValidator interface {
	ValidateCurrentOffset() bool
}

func (o *CurrentOffset) ValidateCurrentOffset() bool {
	return o.valid
}

type CurrentPartNumber struct {
	valid bool
	v     *int32
}

type CurrentPartNumberGetter interface {
	GetCurrentPartNumber() *int32
}

func (o *CurrentPartNumber) GetCurrentPartNumber() *int32 {
	if !o.valid {
		panic("CurrentPartNumber value is not valid")
	}
	return o.v
}

type CurrentPartNumberSetter interface {
	SetCurrentPartNumber(*int32)
}

func (o *CurrentPartNumber) SetCurrentPartNumber(v *int32) {
	o.v = v
	o.valid = true
}

type CurrentPartNumberValidator interface {
	ValidateCurrentPartNumber() bool
}

func (o *CurrentPartNumber) ValidateCurrentPartNumber() bool {
	return o.valid
}

type DestinationService struct {
	valid bool
	v     storage.Servicer
}

type DestinationServiceGetter interface {
	GetDestinationService() storage.Servicer
}

func (o *DestinationService) GetDestinationService() storage.Servicer {
	if !o.valid {
		panic("DestinationService value is not valid")
	}
	return o.v
}

type DestinationServiceSetter interface {
	SetDestinationService(storage.Servicer)
}

func (o *DestinationService) SetDestinationService(v storage.Servicer) {
	o.v = v
	o.valid = true
}

type DestinationServiceValidator interface {
	ValidateDestinationService() bool
}

func (o *DestinationService) ValidateDestinationService() bool {
	return o.valid
}

type DestinationStorage struct {
	valid bool
	v     storage.Storager
}

type DestinationStorageGetter interface {
	GetDestinationStorage() storage.Storager
}

func (o *DestinationStorage) GetDestinationStorage() storage.Storager {
	if !o.valid {
		panic("DestinationStorage value is not valid")
	}
	return o.v
}

type DestinationStorageSetter interface {
	SetDestinationStorage(storage.Storager)
}

func (o *DestinationStorage) SetDestinationStorage(v storage.Storager) {
	o.v = v
	o.valid = true
}

type DestinationStorageValidator interface {
	ValidateDestinationStorage() bool
}

func (o *DestinationStorage) ValidateDestinationStorage() bool {
	return o.valid
}

type EnableBenchmark struct {
	valid bool
	v     bool
}

type EnableBenchmarkGetter interface {
	GetEnableBenchmark() bool
}

func (o *EnableBenchmark) GetEnableBenchmark() bool {
	if !o.valid {
		panic("EnableBenchmark value is not valid")
	}
	return o.v
}

type EnableBenchmarkSetter interface {
	SetEnableBenchmark(bool)
}

func (o *EnableBenchmark) SetEnableBenchmark(v bool) {
	o.v = v
	o.valid = true
}

type EnableBenchmarkValidator interface {
	ValidateEnableBenchmark() bool
}

func (o *EnableBenchmark) ValidateEnableBenchmark() bool {
	return o.valid
}

type ExpectSize struct {
	valid bool
	v     int64
}

type ExpectSizeGetter interface {
	GetExpectSize() int64
}

func (o *ExpectSize) GetExpectSize() int64 {
	if !o.valid {
		panic("ExpectSize value is not valid")
	}
	return o.v
}

type ExpectSizeSetter interface {
	SetExpectSize(int64)
}

func (o *ExpectSize) SetExpectSize(v int64) {
	o.v = v
	o.valid = true
}

type ExpectSizeValidator interface {
	ValidateExpectSize() bool
}

func (o *ExpectSize) ValidateExpectSize() bool {
	return o.valid
}

type Expire struct {
	valid bool
	v     int
}

type ExpireGetter interface {
	GetExpire() int
}

func (o *Expire) GetExpire() int {
	if !o.valid {
		panic("Expire value is not valid")
	}
	return o.v
}

type ExpireSetter interface {
	SetExpire(int)
}

func (o *Expire) SetExpire(v int) {
	o.v = v
	o.valid = true
}

type ExpireValidator interface {
	ValidateExpire() bool
}

func (o *Expire) ValidateExpire() bool {
	return o.valid
}

type Fault struct {
	valid bool
	v     error
}

type FaultGetter interface {
	GetFault() error
}

func (o *Fault) GetFault() error {
	if !o.valid {
		panic("Fault value is not valid")
	}
	return o.v
}

type FaultSetter interface {
	SetFault(error)
}

func (o *Fault) SetFault(v error) {
	o.v = v
	o.valid = true
}

type FaultValidator interface {
	ValidateFault() bool
}

func (o *Fault) ValidateFault() bool {
	return o.valid
}

type FlowType struct {
	valid bool
	v     constants.FlowType
}

type FlowTypeGetter interface {
	GetFlowType() constants.FlowType
}

func (o *FlowType) GetFlowType() constants.FlowType {
	if !o.valid {
		panic("FlowType value is not valid")
	}
	return o.v
}

type FlowTypeSetter interface {
	SetFlowType(constants.FlowType)
}

func (o *FlowType) SetFlowType(v constants.FlowType) {
	o.v = v
	o.valid = true
}

type FlowTypeValidator interface {
	ValidateFlowType() bool
}

func (o *FlowType) ValidateFlowType() bool {
	return o.valid
}

type Force struct {
	valid bool
	v     bool
}

type ForceGetter interface {
	GetForce() bool
}

func (o *Force) GetForce() bool {
	if !o.valid {
		panic("Force value is not valid")
	}
	return o.v
}

type ForceSetter interface {
	SetForce(bool)
}

func (o *Force) SetForce(v bool) {
	o.v = v
	o.valid = true
}

type ForceValidator interface {
	ValidateForce() bool
}

func (o *Force) ValidateForce() bool {
	return o.valid
}

type HumanReadable struct {
	valid bool
	v     bool
}

type HumanReadableGetter interface {
	GetHumanReadable() bool
}

func (o *HumanReadable) GetHumanReadable() bool {
	if !o.valid {
		panic("HumanReadable value is not valid")
	}
	return o.v
}

type HumanReadableSetter interface {
	SetHumanReadable(bool)
}

func (o *HumanReadable) SetHumanReadable(v bool) {
	o.v = v
	o.valid = true
}

type HumanReadableValidator interface {
	ValidateHumanReadable() bool
}

func (o *HumanReadable) ValidateHumanReadable() bool {
	return o.valid
}

type ID struct {
	valid bool
	v     string
}

type IDGetter interface {
	GetID() string
}

func (o *ID) GetID() string {
	if !o.valid {
		panic("ID value is not valid")
	}
	return o.v
}

type IDSetter interface {
	SetID(string)
}

func (o *ID) SetID(v string) {
	o.v = v
	o.valid = true
}

type IDValidator interface {
	ValidateID() bool
}

func (o *ID) ValidateID() bool {
	return o.valid
}

type Key struct {
	valid bool
	v     string
}

type KeyGetter interface {
	GetKey() string
}

func (o *Key) GetKey() string {
	if !o.valid {
		panic("Key value is not valid")
	}
	return o.v
}

type KeySetter interface {
	SetKey(string)
}

func (o *Key) SetKey(v string) {
	o.v = v
	o.valid = true
}

type KeyValidator interface {
	ValidateKey() bool
}

func (o *Key) ValidateKey() bool {
	return o.valid
}

type KeyType struct {
	valid bool
	v     constants.KeyType
}

type KeyTypeGetter interface {
	GetKeyType() constants.KeyType
}

func (o *KeyType) GetKeyType() constants.KeyType {
	if !o.valid {
		panic("KeyType value is not valid")
	}
	return o.v
}

type KeyTypeSetter interface {
	SetKeyType(constants.KeyType)
}

func (o *KeyType) SetKeyType(v constants.KeyType) {
	o.v = v
	o.valid = true
}

type KeyTypeValidator interface {
	ValidateKeyType() bool
}

func (o *KeyType) ValidateKeyType() bool {
	return o.valid
}

type ListType struct {
	valid bool
	v     constants.ListType
}

type ListTypeGetter interface {
	GetListType() constants.ListType
}

func (o *ListType) GetListType() constants.ListType {
	if !o.valid {
		panic("ListType value is not valid")
	}
	return o.v
}

type ListTypeSetter interface {
	SetListType(constants.ListType)
}

func (o *ListType) SetListType(v constants.ListType) {
	o.v = v
	o.valid = true
}

type ListTypeValidator interface {
	ValidateListType() bool
}

func (o *ListType) ValidateListType() bool {
	return o.valid
}

type LongFormat struct {
	valid bool
	v     bool
}

type LongFormatGetter interface {
	GetLongFormat() bool
}

func (o *LongFormat) GetLongFormat() bool {
	if !o.valid {
		panic("LongFormat value is not valid")
	}
	return o.v
}

type LongFormatSetter interface {
	SetLongFormat(bool)
}

func (o *LongFormat) SetLongFormat(v bool) {
	o.v = v
	o.valid = true
}

type LongFormatValidator interface {
	ValidateLongFormat() bool
}

func (o *LongFormat) ValidateLongFormat() bool {
	return o.valid
}

type MD5Sum struct {
	valid bool
	v     []byte
}

type MD5SumGetter interface {
	GetMD5Sum() []byte
}

func (o *MD5Sum) GetMD5Sum() []byte {
	if !o.valid {
		panic("MD5Sum value is not valid")
	}
	return o.v
}

type MD5SumSetter interface {
	SetMD5Sum([]byte)
}

func (o *MD5Sum) SetMD5Sum(v []byte) {
	o.v = v
	o.valid = true
}

type MD5SumValidator interface {
	ValidateMD5Sum() bool
}

func (o *MD5Sum) ValidateMD5Sum() bool {
	return o.valid
}

type Name struct {
	valid bool
	v     string
}

type NameGetter interface {
	GetName() string
}

func (o *Name) GetName() string {
	if !o.valid {
		panic("Name value is not valid")
	}
	return o.v
}

type NameSetter interface {
	SetName(string)
}

func (o *Name) SetName(v string) {
	o.v = v
	o.valid = true
}

type NameValidator interface {
	ValidateName() bool
}

func (o *Name) ValidateName() bool {
	return o.valid
}

type Object struct {
	valid bool
	v     *types.Object
}

type ObjectGetter interface {
	GetObject() *types.Object
}

func (o *Object) GetObject() *types.Object {
	if !o.valid {
		panic("Object value is not valid")
	}
	return o.v
}

type ObjectSetter interface {
	SetObject(*types.Object)
}

func (o *Object) SetObject(v *types.Object) {
	o.v = v
	o.valid = true
}

type ObjectValidator interface {
	ValidateObject() bool
}

func (o *Object) ValidateObject() bool {
	return o.valid
}

type ObjectChannel struct {
	valid bool
	v     chan *types.Object
}

type ObjectChannelGetter interface {
	GetObjectChannel() chan *types.Object
}

func (o *ObjectChannel) GetObjectChannel() chan *types.Object {
	if !o.valid {
		panic("ObjectChannel value is not valid")
	}
	return o.v
}

type ObjectChannelSetter interface {
	SetObjectChannel(chan *types.Object)
}

func (o *ObjectChannel) SetObjectChannel(v chan *types.Object) {
	o.v = v
	o.valid = true
}

type ObjectChannelValidator interface {
	ValidateObjectChannel() bool
}

func (o *ObjectChannel) ValidateObjectChannel() bool {
	return o.valid
}

type ObjectLongList struct {
	valid bool
	v     *[][]string
}

type ObjectLongListGetter interface {
	GetObjectLongList() *[][]string
}

func (o *ObjectLongList) GetObjectLongList() *[][]string {
	if !o.valid {
		panic("ObjectLongList value is not valid")
	}
	return o.v
}

type ObjectLongListSetter interface {
	SetObjectLongList(*[][]string)
}

func (o *ObjectLongList) SetObjectLongList(v *[][]string) {
	o.v = v
	o.valid = true
}

type ObjectLongListValidator interface {
	ValidateObjectLongList() bool
}

func (o *ObjectLongList) ValidateObjectLongList() bool {
	return o.valid
}

type Offset struct {
	valid bool
	v     int64
}

type OffsetGetter interface {
	GetOffset() int64
}

func (o *Offset) GetOffset() int64 {
	if !o.valid {
		panic("Offset value is not valid")
	}
	return o.v
}

type OffsetSetter interface {
	SetOffset(int64)
}

func (o *Offset) SetOffset(v int64) {
	o.v = v
	o.valid = true
}

type OffsetValidator interface {
	ValidateOffset() bool
}

func (o *Offset) ValidateOffset() bool {
	return o.valid
}

type PartNumber struct {
	valid bool
	v     int
}

type PartNumberGetter interface {
	GetPartNumber() int
}

func (o *PartNumber) GetPartNumber() int {
	if !o.valid {
		panic("PartNumber value is not valid")
	}
	return o.v
}

type PartNumberSetter interface {
	SetPartNumber(int)
}

func (o *PartNumber) SetPartNumber(v int) {
	o.v = v
	o.valid = true
}

type PartNumberValidator interface {
	ValidatePartNumber() bool
}

func (o *PartNumber) ValidatePartNumber() bool {
	return o.valid
}

type PartSize struct {
	valid bool
	v     int64
}

type PartSizeGetter interface {
	GetPartSize() int64
}

func (o *PartSize) GetPartSize() int64 {
	if !o.valid {
		panic("PartSize value is not valid")
	}
	return o.v
}

type PartSizeSetter interface {
	SetPartSize(int64)
}

func (o *PartSize) SetPartSize(v int64) {
	o.v = v
	o.valid = true
}

type PartSizeValidator interface {
	ValidatePartSize() bool
}

func (o *PartSize) ValidatePartSize() bool {
	return o.valid
}

type Path struct {
	valid bool
	v     string
}

type PathGetter interface {
	GetPath() string
}

func (o *Path) GetPath() string {
	if !o.valid {
		panic("Path value is not valid")
	}
	return o.v
}

type PathSetter interface {
	SetPath(string)
}

func (o *Path) SetPath(v string) {
	o.v = v
	o.valid = true
}

type PathValidator interface {
	ValidatePath() bool
}

func (o *Path) ValidatePath() bool {
	return o.valid
}

type PathType struct {
	valid bool
	v     constants.PathType
}

type PathTypeGetter interface {
	GetPathType() constants.PathType
}

func (o *PathType) GetPathType() constants.PathType {
	if !o.valid {
		panic("PathType value is not valid")
	}
	return o.v
}

type PathTypeSetter interface {
	SetPathType(constants.PathType)
}

func (o *PathType) SetPathType(v constants.PathType) {
	o.v = v
	o.valid = true
}

type PathTypeValidator interface {
	ValidatePathType() bool
}

func (o *PathType) ValidatePathType() bool {
	return o.valid
}

type Pool struct {
	valid bool
	v     *navvy.Pool
}

type PoolGetter interface {
	GetPool() *navvy.Pool
}

func (o *Pool) GetPool() *navvy.Pool {
	if !o.valid {
		panic("Pool value is not valid")
	}
	return o.v
}

type PoolSetter interface {
	SetPool(*navvy.Pool)
}

func (o *Pool) SetPool(v *navvy.Pool) {
	o.v = v
	o.valid = true
}

type PoolValidator interface {
	ValidatePool() bool
}

func (o *Pool) ValidatePool() bool {
	return o.valid
}

type ReadableSize struct {
	valid bool
	v     string
}

type ReadableSizeGetter interface {
	GetReadableSize() string
}

func (o *ReadableSize) GetReadableSize() string {
	if !o.valid {
		panic("ReadableSize value is not valid")
	}
	return o.v
}

type ReadableSizeSetter interface {
	SetReadableSize(string)
}

func (o *ReadableSize) SetReadableSize(v string) {
	o.v = v
	o.valid = true
}

type ReadableSizeValidator interface {
	ValidateReadableSize() bool
}

func (o *ReadableSize) ValidateReadableSize() bool {
	return o.valid
}

type Recursive struct {
	valid bool
	v     bool
}

type RecursiveGetter interface {
	GetRecursive() bool
}

func (o *Recursive) GetRecursive() bool {
	if !o.valid {
		panic("Recursive value is not valid")
	}
	return o.v
}

type RecursiveSetter interface {
	SetRecursive(bool)
}

func (o *Recursive) SetRecursive(v bool) {
	o.v = v
	o.valid = true
}

type RecursiveValidator interface {
	ValidateRecursive() bool
}

func (o *Recursive) ValidateRecursive() bool {
	return o.valid
}

type Scheduler struct {
	valid bool
	v     scheduler
}

type SchedulerGetter interface {
	GetScheduler() scheduler
}

func (o *Scheduler) GetScheduler() scheduler {
	if !o.valid {
		panic("Scheduler value is not valid")
	}
	return o.v
}

type SchedulerSetter interface {
	SetScheduler(scheduler)
}

func (o *Scheduler) SetScheduler(v scheduler) {
	o.v = v
	o.valid = true
}

type SchedulerValidator interface {
	ValidateScheduler() bool
}

func (o *Scheduler) ValidateScheduler() bool {
	return o.valid
}

type SegmentID struct {
	valid bool
	v     string
}

type SegmentIDGetter interface {
	GetSegmentID() string
}

func (o *SegmentID) GetSegmentID() string {
	if !o.valid {
		panic("SegmentID value is not valid")
	}
	return o.v
}

type SegmentIDSetter interface {
	SetSegmentID(string)
}

func (o *SegmentID) SetSegmentID(v string) {
	o.v = v
	o.valid = true
}

type SegmentIDValidator interface {
	ValidateSegmentID() bool
}

func (o *SegmentID) ValidateSegmentID() bool {
	return o.valid
}

type Size struct {
	valid bool
	v     int64
}

type SizeGetter interface {
	GetSize() int64
}

func (o *Size) GetSize() int64 {
	if !o.valid {
		panic("Size value is not valid")
	}
	return o.v
}

type SizeSetter interface {
	SetSize(int64)
}

func (o *Size) SetSize(v int64) {
	o.v = v
	o.valid = true
}

type SizeValidator interface {
	ValidateSize() bool
}

func (o *Size) ValidateSize() bool {
	return o.valid
}

type SourceStorage struct {
	valid bool
	v     storage.Storager
}

type SourceStorageGetter interface {
	GetSourceStorage() storage.Storager
}

func (o *SourceStorage) GetSourceStorage() storage.Storager {
	if !o.valid {
		panic("SourceStorage value is not valid")
	}
	return o.v
}

type SourceStorageSetter interface {
	SetSourceStorage(storage.Storager)
}

func (o *SourceStorage) SetSourceStorage(v storage.Storager) {
	o.v = v
	o.valid = true
}

type SourceStorageValidator interface {
	ValidateSourceStorage() bool
}

func (o *SourceStorage) ValidateSourceStorage() bool {
	return o.valid
}

type Stream struct {
	valid bool
	v     io.Reader
}

type StreamGetter interface {
	GetStream() io.Reader
}

func (o *Stream) GetStream() io.Reader {
	if !o.valid {
		panic("Stream value is not valid")
	}
	return o.v
}

type StreamSetter interface {
	SetStream(io.Reader)
}

func (o *Stream) SetStream(v io.Reader) {
	o.v = v
	o.valid = true
}

type StreamValidator interface {
	ValidateStream() bool
}

func (o *Stream) ValidateStream() bool {
	return o.valid
}

type TotalSize struct {
	valid bool
	v     int64
}

type TotalSizeGetter interface {
	GetTotalSize() int64
}

func (o *TotalSize) GetTotalSize() int64 {
	if !o.valid {
		panic("TotalSize value is not valid")
	}
	return o.v
}

type TotalSizeSetter interface {
	SetTotalSize(int64)
}

func (o *TotalSize) SetTotalSize(v int64) {
	o.v = v
	o.valid = true
}

type TotalSizeValidator interface {
	ValidateTotalSize() bool
}

func (o *TotalSize) ValidateTotalSize() bool {
	return o.valid
}

type URL struct {
	valid bool
	v     string
}

type URLGetter interface {
	GetURL() string
}

func (o *URL) GetURL() string {
	if !o.valid {
		panic("URL value is not valid")
	}
	return o.v
}

type URLSetter interface {
	SetURL(string)
}

func (o *URL) SetURL(v string) {
	o.v = v
	o.valid = true
}

type URLValidator interface {
	ValidateURL() bool
}

func (o *URL) ValidateURL() bool {
	return o.valid
}

type Zone struct {
	valid bool
	v     string
}

type ZoneGetter interface {
	GetZone() string
}

func (o *Zone) GetZone() string {
	if !o.valid {
		panic("Zone value is not valid")
	}
	return o.v
}

type ZoneSetter interface {
	SetZone(string)
}

func (o *Zone) SetZone(v string) {
	o.v = v
	o.valid = true
}

type ZoneValidator interface {
	ValidateZone() bool
}

func (o *Zone) ValidateZone() bool {
	return o.valid
}
