package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/storage"
)

type BucketTestSuite struct {
	suite.Suite
}

func (suite BucketTestSuite) SetupTest() {
	stor = storage.NewMockObjectStorage()
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
		// Package handler
		input := BucketHandler{}
		assert.Equal(suite.T(), c.expected, input.WithZone(c.zone).WithRemote(c.name).MakeBucket())
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
		// Package handler
		input := BucketHandler{}
		assert.Equal(suite.T(), c.expected1, input.WithZone(c.zone).ListBuckets(), c.zone)
		buckets, _ := stor.ListBuckets(c.zone)
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
		// Package handler
		input := BucketHandler{}
		assert.Equal(suite.T(), c.expected1, input.WithRemote(c.name).RemoveBucket(), c.name)
		buckets, _ := stor.ListBuckets("")
		assert.Equal(suite.T(), c.expected2, len(buckets), c.name)
	}
}

func TestBucketTestSuite(t *testing.T) {
	suite.Run(t, new(BucketTestSuite))
}
