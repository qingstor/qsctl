package action

// BaseHandler contains all flags from cmd
type BaseHandler struct {
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
