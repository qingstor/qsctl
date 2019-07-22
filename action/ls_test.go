package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
)

const (
	objNum    int    = 20
	objPrefix string = "obj"
)

type LsTestSuite struct {
	suite.Suite
}

func (suite LsTestSuite) SetupTest() {
	contexts.Bench = true
	s := storage.NewMockObjectStorage()
	s.AddMockObjects(objPrefix, objNum)
	contexts.Storage = s
}

func (suite LsTestSuite) TestListObjects() {
	cases := []struct {
		remote        string
		key           string
		humanReadable bool
		longFormat    bool
		recursive     bool
		expected      int
		err           error
	}{
		// Add / as postfix to simulate the non-recursive situation
		{fmt.Sprintf("qs://%s//%s", storage.MockPek3a, objPrefix+"/"),
			objPrefix + "/", true, false, false,
			objNum + 1, nil},
		{fmt.Sprintf("qs://%s//%s", storage.MockPek3a, objPrefix),
			objPrefix, true, false, true,
			2*objNum + 1, nil},
		{fmt.Sprintf("qs://%s//%s", storage.MockPek3a, objPrefix),
			objPrefix, false, true, true,
			2*objNum + 1, nil},
		{fmt.Sprintf("qs://%s//%s", storage.MockPek3a, objPrefix),
			objPrefix, true, true, true,
			2*objNum + 1, nil},
	}
	contexts.LongFormat = false
	for _, c := range cases {
		contexts.HumanReadable = c.humanReadable
		contexts.LongFormat = c.longFormat
		contexts.Recursive = c.recursive
		delimiter := "/"
		if c.recursive {
			delimiter = ""
		}
		assert.Equal(suite.T(), ListObjects(c.remote), c.err)
		om, err := contexts.Storage.ListObjects(c.key, delimiter, nil)
		assert.Equal(suite.T(), err, c.err)
		assert.Equal(suite.T(), len(om), c.expected)
	}

}

func TestLsTestSuite(t *testing.T) {
	suite.Run(t, new(LsTestSuite))
}
