package action

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
)

type BucketTestSuite struct {
	suite.Suite
}

func (suite BucketTestSuite) SetupTest() {
	contexts.Storage = storage.NewMockObjectStorage()
}

func (suite BucketTestSuite) TestMakeBucket() {
	cases := []struct {
		zone     string
		name     string
		expected error
	}{
		{storage.MockZoneAlpha, "new-pek3a-bucket", nil},
		{storage.MockZoneAlpha, "new-pek3a-bucket", constants.ErrorBucketAlreadyExists},
		{storage.MockZoneBeta, "qs://new-pek3b-bucket", nil},
		{storage.MockZoneBeta, "qs://", constants.ErrorQsPathInvalid},
	}
	for _, c := range cases {
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, c.zone)
		ctx = contexts.SetContext(ctx, "remote", c.name)

		assert.Equal(suite.T(), c.expected, MakeBucket(ctx))
	}
}

func (suite BucketTestSuite) TestListBuckets() {
	cases := []struct {
		zone      string
		expected1 error
		expected2 int
	}{
		{storage.MockZoneAlpha, nil, 1},
		{storage.MockZoneBeta, nil, 1},
		{"", nil, 2},
	}
	for _, c := range cases {
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, c.zone)

		assert.Equal(suite.T(), c.expected1, ListBuckets(ctx), c.zone)
		buckets, _ := contexts.Storage.ListBuckets(contexts.FromContext(ctx, constants.ZoneFlag).(string))
		assert.Equal(suite.T(), c.expected2, len(buckets), c.zone)
	}
}

func (suite BucketTestSuite) TestRemoveBucket() {
	cases := []struct {
		name      string
		expected1 error
		expected2 int
	}{
		{storage.MockZoneAlpha, nil, 1},
		{storage.MockZoneBeta, nil, 0},
		{"qs://", constants.ErrorQsPathInvalid, 0},
	}
	for _, c := range cases {
		// Package context
		var ctx context.Context
		ctx = contexts.NewMockCmdContext()
		ctx = contexts.SetContext(ctx, constants.ZoneFlag, "")
		ctx = contexts.SetContext(ctx, "remote", c.name)

		assert.Equal(suite.T(), c.expected1, RemoveBucket(ctx), c.name)
		buckets, _ := contexts.Storage.ListBuckets("")
		assert.Equal(suite.T(), c.expected2, len(buckets), c.name)
	}
}

func TestBucketTestSuite(t *testing.T) {
	suite.Run(t, new(BucketTestSuite))
}
