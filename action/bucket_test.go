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
		{storage.MockPek3a, "new-pek3a-bucket", nil},
		{storage.MockPek3a, "new-pek3a-bucket", constants.ErrorBucketAlreadyExists},
		{storage.MockPek3b, "qs://new-pek3b-bucket", nil},
		{storage.MockPek3b, "qs://", constants.ErrorQsPathInvalid},
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
		{storage.MockPek3a, nil, 1},
		{storage.MockSh1a, nil, 1},
		{storage.MockPek3b, nil, 1},
		{storage.MockGd2, nil, 1},
		{"", nil, 4},
	}
	for _, c := range cases {
		assert.Equal(suite.T(), c.expected1, ListBuckets(c.zone))
		buckets, _ := contexts.Storage.ListBuckets(c.zone)
		assert.Equal(suite.T(), c.expected2, len(buckets))
	}
}

func (suite BucketTestSuite) TestRemoveBucket() {
	cases := []struct {
		name      string
		expected1 error
		expected2 int
	}{
		{storage.MockPek3a, nil, 3},
		{storage.MockPek3b, nil, 2},
		{"qs://", constants.ErrorQsPathInvalid, 2},
	}
	contexts.Zone = ""
	for _, c := range cases {
		assert.Equal(suite.T(), c.expected1, RemoveBucket(c.name))
		buckets, _ := contexts.Storage.ListBuckets("")
		assert.Equal(suite.T(), c.expected2, len(buckets))
	}
}

func TestBucketTestSuite(t *testing.T) {
	suite.Run(t, new(BucketTestSuite))
}
