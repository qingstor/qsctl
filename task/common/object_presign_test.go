package common

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/Xuanwo/navvy"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestObjectPresignTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key, bucketName, errKey := uuid.New().String(), uuid.New().String(), "presign-error"
	presignErr := errors.New(errKey)

	store := mock.NewMockStorager(ctrl)

	pool := navvy.NewPool(10)

	cases := []struct {
		name       string
		bucketName string
		key        string
		expire     int
		err        error
	}{
		{"ok", bucketName, key, 100, nil},
		{"error", bucketName, key, -100, presignErr},
	}

	for _, ca := range cases {
		now := time.Now().Unix()
		store.EXPECT().Reach(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, p *types.Pair) (string, error) {
			assert.Equal(t, key, inputPath)
			return strconv.FormatInt(now+int64(ca.expire), 10), ca.err
		}).Times(1)

		x := &mockObjectPresignTask{}
		x.SetPool(pool)
		x.SetDestinationStorage(store)
		x.SetBucketName(ca.bucketName)
		x.SetDestinationPath(ca.key)
		x.SetExpire(ca.expire)

		task := NewObjectPresignTask(x)
		task.Run()
		pool.Wait()

		if ca.err != nil {
			assert.Equal(t, x.ValidateFault(), true)
			assert.Error(t, x.GetFault())
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.err))
			continue
		}
		v, e := strconv.ParseInt(x.GetURL(), 10, 64)
		t.Log(v)
		if e != nil {
			t.Fatal(e)
		}
		assert.Equal(t, false, x.ValidateFault())
		assert.Equal(t, now+int64(ca.expire), v)
	}
}
