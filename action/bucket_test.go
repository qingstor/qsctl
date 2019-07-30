package action

import (
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
	contexts.Bench = true
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
		contexts.Zone = c.zone
		assert.Equal(suite.T(), c.expected, MakeBucket(c.name))
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
		assert.Equal(suite.T(), c.expected1, ListBuckets(c.zone), c.zone)
		buckets, _ := contexts.Storage.ListBuckets(c.zone)
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
	contexts.Zone = ""
	for _, c := range cases {
		assert.Equal(suite.T(), c.expected1, RemoveBucket(c.name), c.name)
		buckets, _ := contexts.Storage.ListBuckets("")
		assert.Equal(suite.T(), c.expected2, len(buckets), c.name)
	}
}

func TestBucketTestSuite(t *testing.T) {
	suite.Run(t, new(BucketTestSuite))
}
