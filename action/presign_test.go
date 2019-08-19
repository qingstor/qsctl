package action

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yunify/qsctl/v2/constants"
	"github.com/yunify/qsctl/v2/contexts"
	"github.com/yunify/qsctl/v2/storage"
)

type PresignTestSuite struct {
	suite.Suite
}

func (suite PresignTestSuite) SetupTest() {
	contexts.Storage = storage.NewMockObjectStorage()
}

func (suite PresignTestSuite) TestPresign() {
	cases := []struct {
		expire int
		remote string
		zone   string
		err    error
	}{
		{10, fmt.Sprintf("qs://%s/%s", storage.MockZoneAlpha, storage.Mock0BObject),
			storage.MockZoneAlpha, nil},
		{10, fmt.Sprintf("qs://%s/%s", storage.MockZoneBeta, storage.Mock0BObject),
			storage.MockZoneBeta, nil},
		{10, fmt.Sprintf("qs://%s/%s", storage.MockZoneBeta, objPrefix),
			storage.MockZoneBeta, constants.ErrorQsPathNotFound},
		{10, "", storage.MockZoneBeta, constants.ErrorQsPathInvalid},
		{10, fmt.Sprintf("qs://%s", storage.MockZoneBeta),
			storage.MockZoneBeta, constants.ErrorQsPathObjectKeyRequired},
		{-1, fmt.Sprintf("qs://%s/%s", storage.MockZoneBeta, storage.Mock0BObject),
			storage.MockZoneBeta, constants.ErrorTestError},
	}

	for _, c := range cases {
		input := PresignHandler{}
		err := input.WithExpire(c.expire).WithRemote(c.remote).WithZone(c.zone).Presign()
		assert.Equal(suite.T(), c.err, err)
	}
}

func TestPresignTestSuite(t *testing.T) {
	suite.Run(t, new(PresignTestSuite))
}
