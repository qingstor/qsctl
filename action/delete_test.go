package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
)

type DeleteTestSuite struct {
	suite.Suite
}

func (suite DeleteTestSuite) SetupTest() {
	contexts.Storage = storage.NewMockObjectStorage()
}

func (suite DeleteTestSuite) TestDelete() {
	zone := ""
	cases := []struct {
		msg       string
		recursive bool
		remote    string
		objCount  int
		err       error
	}{
		// we got 5 * objNum (in for-loop subdirectories) + 2 (dir: obj/ & obj/obj/) + 5 (0B/MB/GB/TB/Forbidden object) mock objects
		{"fail get QS path", false, storage.MockGBObject, 5*objNum + 2 + 5, constants.ErrorQsPathInvalid},
		{"fail get obj key", false, "qs://bucket/", 5*objNum + 2 + 5, constants.ErrorQsPathObjectKeyRequired},
		{"fail head obj", false, "qs://bucket/not-exist", 5*objNum + 2 + 5, constants.ErrorQsPathNotFound},
		{"fail delete directory", false, "qs://bucket/" + objPrefix + "/", 5*objNum + 2 + 5, constants.ErrorRecursiveRequired},
		{"delete directory", true, "qs://bucket/" + objPrefix + "/", objNum + 5, nil},
		{"fail delete obj", false, "qs://bucket/" + storage.MockForbiddenObject, 5*objNum + 2 + 5, constants.ErrorQsPathAccessForbidden},
		{"delete GB object", false, "qs://bucket/" + storage.MockGBObject, 5*objNum + 2 + 4, nil},
	}

	for _, c := range cases {
		input := &DeleteHandler{}
		s := contexts.Storage.(*storage.MockObjectStorage)
		s.ResetMockObjects(objPrefix, objNum)
		assert.Equal(suite.T(), c.err, input.WithRecursive(c.recursive).WithRemote(c.remote).WithZone(zone).Delete(), c.msg)
		assert.Equal(suite.T(), c.objCount, s.GetObjCount(), c.msg)
	}
}

func TestDeleteTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteTestSuite))
}
