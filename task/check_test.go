package task

import (
	"testing"
	"time"

	typ "github.com/Xuanwo/storage/types"
	"github.com/stretchr/testify/assert"
)

func TestIsDestinationObjectExistTask_run(t *testing.T) {
	t.Run("destination object not exist", func(t *testing.T) {
		task := IsDestinationObjectExistTask{}
		task.SetDestinationObject(nil)

		task.run()

		assert.Equal(t, false, task.GetResult())
	})

	t.Run("destination object exists", func(t *testing.T) {
		task := IsDestinationObjectExistTask{}
		task.SetDestinationObject(&typ.Object{})

		task.run()

		assert.Equal(t, true, task.GetResult())
	})
}

func TestIsSizeEqualTask_run(t *testing.T) {
	t.Run("size equal", func(t *testing.T) {
		task := IsSizeEqualTask{}
		task.SetSourceObject(&typ.Object{Size: 111})
		task.SetDestinationObject(&typ.Object{Size: 111})

		task.run()

		assert.Equal(t, true, task.GetResult())
	})

	t.Run("size not equal", func(t *testing.T) {
		task := IsSizeEqualTask{}
		task.SetSourceObject(&typ.Object{Size: 222})
		task.SetDestinationObject(&typ.Object{Size: 111})

		task.run()

		assert.Equal(t, false, task.GetResult())
	})
}

func TestIsUpdateAtGreaterTask_run(t *testing.T) {
	t.Run("updated at greater", func(t *testing.T) {
		task := IsUpdateAtGreaterTask{}
		task.SetSourceObject(&typ.Object{UpdatedAt: time.Now().Add(time.Hour)})
		task.SetDestinationObject(&typ.Object{UpdatedAt: time.Now()})

		task.run()

		assert.Equal(t, true, task.GetResult())
	})

	t.Run("updated at not greater", func(t *testing.T) {
		task := IsUpdateAtGreaterTask{}
		task.SetSourceObject(&typ.Object{UpdatedAt: time.Now()})
		task.SetDestinationObject(&typ.Object{UpdatedAt: time.Now().Add(time.Hour)})

		task.run()

		assert.Equal(t, false, task.GetResult())
	})
}
