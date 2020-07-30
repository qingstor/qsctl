package shellutils

import (
	"github.com/Xuanwo/storage"
	"github.com/qingstor/noah/task"
	log "github.com/sirupsen/logrus"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/utils"
)

var bucketList = make([]string, 0, 10)

// InitBucketList init bucket list as cache
func InitBucketList() {
	rootTask := taskutils.NewAtServiceTask(10)
	err := utils.ParseAtServiceInput(rootTask)
	if err != nil {
		log.Errorf("get service failed: [%v]", err)
		return
	}

	t := task.NewListStorage(rootTask)
	t.SetZone("")
	t.SetStoragerFunc(func(stor storage.Storager) {
		sm, _ := stor.Metadata()
		bucketList = append(bucketList, sm.Name)
	})
	t.Run()
	return
}

// GetBucketList get list from cache
func GetBucketList() []string {
	return bucketList
}

// RemoveBucketFromList remove bucket from cache
func RemoveBucketFromList(bucket string) {
	if len(bucketList) == 0 {
		return
	}
	for i, b := range bucketList {
		if b == bucket {
			bucketList = append(bucketList[:i], bucketList[i+1:]...)
			break
		}
	}
	return
}

// AddBucketIntoList add bucket into cache
func AddBucketIntoList(bucket string) {
	bucketList = append(bucketList, bucket)
}
