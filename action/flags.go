package action

// FlagHandler contains all flags from cmd
type FlagHandler struct {
	Bench                bool   `json:"bench"`
	ExpectSize           int64  `json:"expect_size"`
	Format               string `json:"format"`
	HumanReadable        bool   `json:"human_readable"`
	LongFormat           bool   `json:"long_format"`
	MaximumMemoryContent int64  `json:"maximum_memory_content"`
	Recursive            bool   `json:"recursive"`
	Reverse              bool   `json:"reverse"`
	Zone                 string `json:"zone"`
}

// checkNil returns a non-nil pointer of FlagHandler for embed-in use in action handler
func (fh *FlagHandler) checkNil() *FlagHandler {
	if fh == nil {
		return &FlagHandler{}
	}
	return fh
}

// WithBench sets flag handler with given bench flag
func (fh *FlagHandler) WithBench(b bool) *FlagHandler {
	return withBench(fh, b)
}

func withBench(fh *FlagHandler, b bool) *FlagHandler {
	fh = fh.checkNil()
	fh.Bench = b
	return fh
}

// WithExpectSize sets flag handler with given expect_size flag
func (fh *FlagHandler) WithExpectSize(size int64) *FlagHandler {
	return withExpectSize(fh, size)
}

func withExpectSize(fh *FlagHandler, size int64) *FlagHandler {
	fh = fh.checkNil()
	fh.ExpectSize = size
	return fh
}

// WithFormat sets flag handler with given format flag
func (fh *FlagHandler) WithFormat(f string) *FlagHandler {
	return withFormat(fh, f)
}

func withFormat(fh *FlagHandler, f string) *FlagHandler {
	fh = fh.checkNil()
	fh.Format = f
	return fh
}

// WithHumanReadable sets flag handler with given human_readable flag
func (fh *FlagHandler) WithHumanReadable(h bool) *FlagHandler {
	return withHumanReadable(fh, h)
}

func withHumanReadable(fh *FlagHandler, h bool) *FlagHandler {
	fh = fh.checkNil()
	fh.HumanReadable = h
	return fh
}

// WithLongFormat sets flag handler with given long_format flag
func (fh *FlagHandler) WithLongFormat(l bool) *FlagHandler {
	return withLongFormat(fh, l)
}

func withLongFormat(fh *FlagHandler, l bool) *FlagHandler {
	fh = fh.checkNil()
	fh.LongFormat = l
	return fh
}

// WithMaximumMemory sets flag handler with given maximum_memory_content flag
func (fh *FlagHandler) WithMaximumMemory(maxMemory int64) *FlagHandler {
	return withMaximumMemory(fh, maxMemory)
}

func withMaximumMemory(fh *FlagHandler, maxMemory int64) *FlagHandler {
	fh = fh.checkNil()
	fh.MaximumMemoryContent = maxMemory
	return fh
}

// WithRecursive sets flag handler with given recursive flag
func (fh *FlagHandler) WithRecursive(r bool) *FlagHandler {
	return withRecursive(fh, r)
}

func withRecursive(fh *FlagHandler, r bool) *FlagHandler {
	fh = fh.checkNil()
	fh.Recursive = r
	return fh
}

// WithReverse sets flag handler with given reverse flag
func (fh *FlagHandler) WithReverse(r bool) *FlagHandler {
	return withReverse(fh, r)
}

func withReverse(fh *FlagHandler, r bool) *FlagHandler {
	fh = fh.checkNil()
	fh.Reverse = r
	return fh
}

// WithZone sets flag handler with given zone flag
func (fh *FlagHandler) WithZone(z string) *FlagHandler {
	return withZone(fh, z)
}

func withZone(fh *FlagHandler, z string) *FlagHandler {
	fh = fh.checkNil()
	fh.Zone = z
	return fh
}

// GetBench gets bench flag from handler
func (fh *FlagHandler) GetBench() bool {
	return fh.checkNil().Bench
}

// GetExpectSize gets expect-size flag from handler
func (fh *FlagHandler) GetExpectSize() int64 {
	return fh.checkNil().ExpectSize
}

// GetFormat gets format flag from handler
func (fh *FlagHandler) GetFormat() string {
	return fh.checkNil().Format
}

// GetHumanReadable gets human_readable flag from handler
func (fh *FlagHandler) GetHumanReadable() bool {
	return fh.checkNil().HumanReadable
}

// GetLongFormat gets long_format flag from handler
func (fh *FlagHandler) GetLongFormat() bool {
	return fh.checkNil().LongFormat
}

// GetMaximumMemory gets maximum_memory_content flag from handler
func (fh *FlagHandler) GetMaximumMemory() int64 {
	return fh.checkNil().MaximumMemoryContent
}

// GetRecursive gets recursive flag from handler
func (fh *FlagHandler) GetRecursive() bool {
	return fh.checkNil().Recursive
}

// GetReverse gets reverse flag from handler
func (fh *FlagHandler) GetReverse() bool {
	return fh.checkNil().Reverse
}

// GetZone gets zone flag from handler
func (fh *FlagHandler) GetZone() string {
	return fh.checkNil().Zone
}
