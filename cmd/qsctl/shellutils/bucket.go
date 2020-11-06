package shellutils

import (
	"context"
	"errors"
	"fmt"
	"sync"

	typ "github.com/aos-dev/go-storage/v2/types"
	"github.com/qingstor/noah/task"

	"github.com/qingstor/qsctl/v2/cmd/qsctl/taskutils"
	"github.com/qingstor/qsctl/v2/utils"
)

var mu = new(sync.Mutex)

var bucketList = make([]string, 0, 10)

// InitBucketList init bucket list as cache
func InitBucketList(ctx context.Context) error {
	rootTask := taskutils.NewAtServiceTask()
	err := utils.ParseAtServiceInput(rootTask)
	if err != nil {
		return fmt.Errorf("get service failed: [%w]", err)
	}

	t := task.NewListStorage(rootTask)
	if err := t.Run(ctx); err != nil {
		return fmt.Errorf("list storage failed: [%w]", err)
	}

	it := t.GetStorageIter()
	for {
		obj, err := it.Next()
		if err != nil {
			if errors.Is(err, typ.IterateDone) {
				break
			}
			return fmt.Errorf("iterate storage failed: [%w]", err)
		}

		sm, _ := obj.Metadata()
		AddBucketIntoList(sm.Name)
	}

	return nil
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
