package action

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
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
	contexts.Storage = storage.NewMockObjectStorage()
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
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.HumanReadableFlag, c.humanReadable)
		ctx = contexts.SetContext(ctx, constants.LongFormatFlag, c.longFormat)
		ctx = contexts.SetContext(ctx, constants.RecursiveFlag, c.recursive)
		ctx = contexts.SetContext(ctx, constants.ReverseFlag, c.reverse)
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, "")
		ctx = contexts.SetContext(ctx, "remote", c.remote)
		ctx = contexts.SetContext(ctx, "prefix", c.key)
		ctx = contexts.SetContext(ctx, "delimiter", delimiter)
		// Reset mock objects after each call
		s := contexts.Storage.(*storage.MockObjectStorage)
		s.ResetMockObjects(objPrefix, objNum)
		assert.Equal(suite.T(), c.err, ListObjects(ctx), k)

		s.ResetMockObjects(objPrefix, objNum)
		oms, err := contexts.Storage.ListObjects(c.key, delimiter, nil)
		assert.Equal(suite.T(), c.err, err, k)
		assert.Equal(suite.T(), c.omsCount, len(oms), k)

		s.ResetMockObjects(objPrefix, objNum)
		root, _ := listObjects(ctx)
		count := root.ChildrenCount()
		assert.Equal(suite.T(), c.childrenCount, count, k)
	}

}

func TestLsTestSuite(t *testing.T) {
	suite.Run(t, new(LsTestSuite))
}
