// Code generated by go generate; DO NOT EDIT.
package common

import (
	"fmt"

	"github.com/Xuanwo/navvy"
	"github.com/google/uuid"

	"github.com/yunify/qsctl/v2/pkg/types"
)

var _ navvy.Pool
var _ types.Pool
var _ = uuid.New()

// bucketCreateTaskRequirement is the requirement for execute BucketCreateTask.
type bucketCreateTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.BucketNameGetter
	types.DestinationServiceGetter
	types.ZoneGetter

	// Mutable value
}

// mockBucketCreateTask is the mock task for BucketCreateTask.
type mockBucketCreateTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BucketName
	types.DestinationService
	types.Zone

	// Mutable value
}

func (t *mockBucketCreateTask) Run() {
	panic("mockBucketCreateTask should not be run.")
}

// BucketCreateTask will send put request to create a bucket.
type BucketCreateTask struct {
	bucketCreateTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
}

// Run implement navvy.Task
func (t *BucketCreateTask) Run() {
	t.run()
}

func (t *BucketCreateTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task BucketCreate failed: {%w}", err))
}

// Wait will wait until BucketCreateTask has been finished
func (t *BucketCreateTask) Wait() {
	t.GetPool().Wait()
}

// bucketDeleteTaskRequirement is the requirement for execute BucketDeleteTask.
type bucketDeleteTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.BucketNameGetter
	types.DestinationServiceGetter

	// Mutable value
}

// mockBucketDeleteTask is the mock task for BucketDeleteTask.
type mockBucketDeleteTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BucketName
	types.DestinationService

	// Mutable value
}

func (t *mockBucketDeleteTask) Run() {
	panic("mockBucketDeleteTask should not be run.")
}

// BucketDeleteTask will send delete request to delete a bucket.
type BucketDeleteTask struct {
	bucketDeleteTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
}

// Run implement navvy.Task
func (t *BucketDeleteTask) Run() {
	t.run()
}

func (t *BucketDeleteTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task BucketDelete failed: {%w}", err))
}

// Wait will wait until BucketDeleteTask has been finished
func (t *BucketDeleteTask) Wait() {
	t.GetPool().Wait()
}

// bucketListTaskRequirement is the requirement for execute BucketListTask.
type bucketListTaskRequirement interface {
	navvy.Task

	// Inherited value
	types.DestinationServiceGetter
	types.ZoneGetter

	// Mutable value
}

// mockBucketListTask is the mock task for BucketListTask.
type mockBucketListTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.DestinationService
	types.Zone

	// Mutable value
}

func (t *mockBucketListTask) Run() {
	panic("mockBucketListTask should not be run.")
}

// BucketListTask will send get request to get bucket list.
type BucketListTask struct {
	bucketListTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
	types.BucketList
}

// Run implement navvy.Task
func (t *BucketListTask) Run() {
	t.run()
}

func (t *BucketListTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task BucketList failed: {%w}", err))
}

// Wait will wait until BucketListTask has been finished
func (t *BucketListTask) Wait() {
	t.GetPool().Wait()
}

// removeBucketForceTaskRequirement is the requirement for execute RemoveBucketForceTask.
type removeBucketForceTaskRequirement interface {
	navvy.Task
	types.PoolGetter

	// Inherited value
	types.BucketNameGetter
	types.DestinationPathGetter
	types.DestinationServiceGetter
	types.DestinationStorageGetter

	// Mutable value
}

// mockRemoveBucketForceTask is the mock task for RemoveBucketForceTask.
type mockRemoveBucketForceTask struct {
	types.Todo
	types.Pool
	types.Fault
	types.ID

	// Inherited value
	types.BucketName
	types.DestinationPath
	types.DestinationService
	types.DestinationStorage

	// Mutable value
}

func (t *mockRemoveBucketForceTask) Run() {
	panic("mockRemoveBucketForceTask should not be run.")
}

// RemoveBucketForceTask will remove a bucket force.
type RemoveBucketForceTask struct {
	removeBucketForceTaskRequirement

	// Predefined runtime value
	types.Fault
	types.ID
	types.Scheduler

	// Runtime value
	types.Done
	types.ObjectChannel
	types.Recursive
}

// Run implement navvy.Task
func (t *RemoveBucketForceTask) Run() {
	t.run()
}

func (t *RemoveBucketForceTask) TriggerFault(err error) {
	t.SetFault(fmt.Errorf("Task RemoveBucketForce failed: {%w}", err))
}

// NewRemoveBucketForceTask will create a RemoveBucketForceTask and fetch inherited data from RemoveBucketTask.
func NewRemoveBucketForceTask(task navvy.Task) navvy.Task {
	t := &RemoveBucketForceTask{
		removeBucketForceTaskRequirement: task.(removeBucketForceTaskRequirement),
	}
	t.SetID(uuid.New().String())
	t.new()
	return t
}
