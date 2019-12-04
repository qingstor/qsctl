package task

import (
	"testing"

	"github.com/Xuanwo/storage"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/fault"
	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestCreateStorageTask_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := mock.NewMockServicer(ctrl)
	storageName := uuid.New().String()
	zone := uuid.New().String()

	task := CreateStorageTask{}
	task.SetFault(fault.New())
	task.SetService(service)
	task.SetStorageName(storageName)
	task.SetZone(zone)

	service.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(name string, pairs ...*types.Pair) (storage.Storager, error) {
		assert.Equal(t, storageName, name)
		assert.Equal(t, zone, pairs[0].Value.(string))
		return nil, nil
	})

	task.run()
	assert.Empty(t, task.GetFault().Error())
}
