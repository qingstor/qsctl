package types

import (
	"github.com/Xuanwo/navvy"
)

type pathScheduleFunc func(navvy.Task) interface {
	navvy.Task

	PathSetter
}

type segmentScheduleFunc func(navvy.Task) interface {
	navvy.Task

	OffsetSetter
	DoneSetter
	SizeGetter
	DoneGetter
}
