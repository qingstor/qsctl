package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FlagsTestSuite struct {
	suite.Suite
}

func (suite FlagsTestSuite) TestNilCheck() {
	var c1 *FlagHandler
	c2 := &FlagHandler{Bench: true}
	cases := []struct {
		fh    *FlagHandler
		bench bool
		msg   string
	}{
		{c1, false, "nil pointer"},
		{c2, true, "not nil"},
	}

	for _, c := range cases {
		res := c.fh.checkNil()
		assert.Equal(suite.T(), res.Bench, c.bench, c.msg)
	}
}

func (suite FlagsTestSuite) TestFlagSet() {
	f := FlagHandler{
		Bench:                true,
		ExpectSize:           100,
		Format:               "format",
		HumanReadable:        true,
		LongFormat:           true,
		MaximumMemoryContent: 1000,
		Recursive:            true,
		Reverse:              true,
		Zone:                 "zone",
	}
	input := &FlagHandler{}
	input.WithBench(f.Bench).
		WithExpectSize(f.ExpectSize).
		WithFormat(f.Format).
		WithHumanReadable(f.HumanReadable).
		WithLongFormat(f.LongFormat).
		WithMaximumMemory(f.MaximumMemoryContent).
		WithRecursive(f.Recursive).
		WithReverse(f.Reverse).
		WithZone(f.Zone)
	assert.Equal(suite.T(), input.Bench, f.Bench, "Bench")
	assert.Equal(suite.T(), input.ExpectSize, f.ExpectSize, "ExpectSize")
	assert.Equal(suite.T(), input.Format, f.Format, "Format")
	assert.Equal(suite.T(), input.HumanReadable, f.HumanReadable, "HumanReadable")
	assert.Equal(suite.T(), input.LongFormat, f.LongFormat, "LongFormat")
	assert.Equal(suite.T(), input.MaximumMemoryContent, f.MaximumMemoryContent, "MaximumMemoryContent")
	assert.Equal(suite.T(), input.Recursive, f.Recursive, "Recursive")
	assert.Equal(suite.T(), input.Reverse, f.Reverse, "Reverse")
	assert.Equal(suite.T(), input.Zone, f.Zone, "Zone")
}

func (suite FlagsTestSuite) TestFlagGet() {
	f := FlagHandler{}
	var input *FlagHandler
	assert.Equal(suite.T(), input.GetBench(), f.Bench, "Bench")
	assert.Equal(suite.T(), input.GetExpectSize(), f.ExpectSize, "ExpectSize")
	assert.Equal(suite.T(), input.GetFormat(), f.Format, "Format")
	assert.Equal(suite.T(), input.GetHumanReadable(), f.HumanReadable, "HumanReadable")
	assert.Equal(suite.T(), input.GetLongFormat(), f.LongFormat, "LongFormat")
	assert.Equal(suite.T(), input.GetMaximumMemory(), f.MaximumMemoryContent, "MaximumMemoryContent")
	assert.Equal(suite.T(), input.GetRecursive(), f.Recursive, "Recursive")
	assert.Equal(suite.T(), input.GetReverse(), f.Reverse, "Reverse")
	assert.Equal(suite.T(), input.GetZone(), f.Zone, "Zone")
}

func TestFlagsTestSuite(t *testing.T) {
	suite.Run(t, new(FlagsTestSuite))
}
