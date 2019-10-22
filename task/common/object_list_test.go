package common

import (
	"errors"
	"testing"

	"github.com/Xuanwo/storage/pkg/iterator"
	"github.com/Xuanwo/storage/types"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/yunify/qsctl/v2/pkg/mock"
)

func TestObjectListTask_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStorager(ctrl)
	key, listErr := uuid.New().String(), errors.New("list-object-err")

	cases := []struct {
		name      string
		key       string
		recursive bool
		fault     bool
		err       error
	}{
		{"non-recursive ok", key, false, false, iterator.ErrDone},
		{"recursive ok", key, true, false, iterator.ErrDone},
		{"non-recursive error not done", key, false, true, listErr},
		{"recursive error not done", key, true, true, listErr},
	}

	for _, ca := range cases {
		x := &mockObjectListTask{}
		x.SetDestinationStorage(store)
		x.SetDestinationPath(ca.key)
		x.SetRecursive(ca.recursive)
		x.SetObjectChannel(make(chan *types.Object))

		go func() {
			// make channel not blocked when set
			for {
				<-x.GetObjectChannel()
			}
		}()

		store.EXPECT().ListDir(gomock.Any(), gomock.Any()).DoAndReturn(func(inputPath string, paris ...*types.Pair) iterator.ObjectIterator {
			assert.Equal(t, inputPath, key)
			if ca.recursive {
				assert.Equal(t, 0, len(paris), ca.name)
			} else {
				assert.Equal(t, "/", paris[0].Value.(string), ca.name)
			}
			count := 3
			return iterator.NewObjectIterator(func(object *[]*types.Object) error {
				*object = make([]*types.Object, 1)
				count--
				if count > 0 {
					return nil
				}
				return ca.err
			})
		})

		task := NewObjectListTask(x)
		task.Run()

		assert.Equal(t, ca.fault, x.ValidateFault(), ca.name)
		if ca.fault {
			assert.Error(t, x.GetFault(), ca.name)
			assert.Equal(t, true, errors.Is(x.GetFault(), ca.err), ca.name)
		}
	}
}
