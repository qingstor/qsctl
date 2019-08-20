package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/storage"
)

const (
	objNum    int    = 5
	objPrefix string = "obj"
)

type LsTestSuite struct {
	suite.Suite
}

func (suite LsTestSuite) SetupTest() {
	stor = storage.NewMockObjectStorage()
}

func (suite LsTestSuite) TestListObjects() {
	cases := []struct {
		remote        string
		key           string
		humanReadable bool
		longFormat    bool
		recursive     bool
		reverse       bool
		omsCount      int
		childrenCount int
		err           error
	}{
		// ls qs://alpha/obj
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix), objPrefix,
			false, false, false, false,
			objNum + 1, objNum + 1, nil},
		// ls qs://alpha/obj -l
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix), objPrefix,
			false, true, false, false,
			objNum + 1, 2*objNum + 2, nil},
		// ls qs://alpha/obj -lh
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix), objPrefix,
			true, true, false, false,
			objNum + 1, 2*objNum + 2, nil},
		// ls qs://alpha/obj -lhRr
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix), objPrefix,
			true, true, true, true,
			objNum + 1, 5*objNum + 2, nil},
		// ls qs://alpha/obj/ -r
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix+"/"), objPrefix + "/",
			false, false, false, true,
			objNum + 2, objNum + 1, nil},
		// ls qs://alpha/obj/ -l
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix+"/"), objPrefix + "/",
			false, true, false, false,
			objNum + 2, 3*objNum + 1, nil},
		// ls qs://alpha/obj/ -lhRr
		{fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, objPrefix+"/"), objPrefix + "/",
			true, true, true, true,
			objNum + 2, 4*objNum + 1, nil},
	}

	delimiter := "/"
	for k, c := range cases {
		// Package handler
		input := &ListHandler{}
		input = input.WithHumanReadable(c.humanReadable).WithLongFormat(c.longFormat).WithRecursive(c.recursive).
			WithReverse(c.reverse).WithRemote(c.remote).WithPrefix(c.key).WithDelimiter(delimiter)
		// Reset mock objects after each call
		s := stor.(*storage.MockObjectStorage)
		s.ResetMockObjects(objPrefix, objNum)
		assert.Equal(suite.T(), c.err, input.ListObjects(), k)

		s.ResetMockObjects(objPrefix, objNum)
		oms, err := stor.ListObjects(c.key, delimiter, nil)
		assert.Equal(suite.T(), c.err, err, k)
		assert.Equal(suite.T(), c.omsCount, len(oms), k)

		s.ResetMockObjects(objPrefix, objNum)
		root, _ := input.listObjects()
		count := root.ChildrenCount()
		assert.Equal(suite.T(), c.childrenCount, count, k)
	}

}

func TestLsTestSuite(t *testing.T) {
	suite.Run(t, new(LsTestSuite))
}
