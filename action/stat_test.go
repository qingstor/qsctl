package action

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
)

type StatTestSuite struct {
	suite.Suite
}

func (suite StatTestSuite) SetupTest() {
	contexts.Storage = storage.NewMockObjectStorage()
}

func (suite StatTestSuite) TestStat() {
	cases := []struct {
		msg   string
		input string
		err   error
	}{
		{"stat object", "qs://bucket/" + storage.Mock0BObject, nil},
	}

	for _, v := range cases {
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.FormatFlag, "")
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, "")
		ctx = contexts.SetContext(ctx, "remote", v.input)

		err := Stat(ctx)
		assert.Equal(suite.T(), v.err, err, v.msg)
	}
}

func (suite StatTestSuite) TestStatWithFormat() {
	cases := []struct {
		msg             string
		input           string
		format          string
		expectedContent string
		err             error
	}{
		{"stat object", "qs://bucket/" + storage.Mock0BObject, "%s", "0\n", nil},
	}

	for _, v := range cases {
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.FormatFlag, v.format)
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, "")
		ctx = contexts.SetContext(ctx, "remote", v.input)

		tempfile, err := ioutil.TempFile("", uuid.New().String())
		if err != nil {
			panic(err)
		}
		defer os.Remove(tempfile.Name())

		os.Stdout = tempfile

		err = Stat(ctx)
		assert.Equal(suite.T(), v.err, err, v.msg)

		_, err = tempfile.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		content, err := ioutil.ReadAll(tempfile)
		if err != nil {
			panic(err)
		}

		assert.Equal(suite.T(), v.expectedContent, string(content), v.msg)
	}
}

func TestStatTestSuite(t *testing.T) {
	suite.Run(t, new(StatTestSuite))
}
