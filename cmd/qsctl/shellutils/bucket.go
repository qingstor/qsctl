package shellutils

import (
	"sync"

	"github.com/Xuanwo/storage"
	"github.com/qingstor/noah/task"
	log "github.com/sirupsen/logrus"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/utils"
)

var mu = new(sync.Mutex)

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
		mu.Lock()
		bucketList = append(bucketList, sm.Name)
		mu.Unlock()
	})
	t.Run()
	return
}

// GetBucketList copy list from cache to avoid data race
func GetBucketList() []string {
	mu.Lock()
	defer mu.Unlock()
	res := make([]string, len(bucketList))
	copy(res, bucketList)
	return res
}

// RemoveBucketFromList remove bucket from cache
func RemoveBucketFromList(bucket string) {
	if len(bucketList) == 0 {
		return
	}
	for i, b := range bucketList {
		if b == bucket {
			mu.Lock()
			bucketList = append(bucketList[:i], bucketList[i+1:]...)
			mu.Unlock()
			break
		}
	}
	return
}

// AddBucketIntoList add bucket into cache
func AddBucketIntoList(bucket string) {
	mu.Lock()
	bucketList = append(bucketList, bucket)
	mu.Unlock()
}
