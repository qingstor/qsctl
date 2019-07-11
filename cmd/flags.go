package cmd

// FlagSet used for registering CtlFlags to a cmd, key is flag name and value is flag info struct
type FlagSet map[string]FlagImp

// CommandToFlagSet used for registering all command-FlagSet, key is command name and value is FlagSet of the command
type CommandToFlagSet map[string]*FlagSet

var cmdToFlagSet = CommandToFlagSet{}

const (
	// all flags' input here
	expectSizeFlag           = "expect-size"
	maximumMemoryContentFlag = "maximum-memory-content"
	zoneFlag                 = "zone"
	recursiveFlag            = "recursive"
)

var (
	// register available flag vars here
	expectSize           string
	maximumMemoryContent string
	zone                 string
	recursive            bool
)

// AddFlag will register a flag with specific flag name, notice that new flag will overwrite old one.
func (fm FlagSet) AddFlag(name string, cf FlagImp) {
	fm[name] = cf
}

// AddFlagSet will register a flag map with its cmd name, notice that new flag map will overwrite old one.
func (cfm CommandToFlagSet) AddFlagSet(cmdName string, fm *FlagSet) {
	cfm[cmdName] = fm
}

// GetRequiredFlags get required flag names for a specific command
func (cfm CommandToFlagSet) GetRequiredFlags(cmdName string) (requiredFlags []FlagImp) {
	flagSet, ok := cfm[cmdName]
	// if the given command has no flags, return
	if !ok {
		return
	}
	for _, ctlFlag := range *flagSet {
		if ctlFlag.GetRequired() {
			requiredFlags = append(requiredFlags, ctlFlag)
		}
	}
	return
}

// FlagImp implements methods for command flag
type FlagImp interface {
	GetName() string
	GetRequired() bool
	SetRequired() FlagImp
	SetPattern(string) FlagImp
	CheckRequired(string) bool
}

// BaseCtlFlag used for base flag infos
type BaseCtlFlag struct {
	Name      string
	ShortHand string
	Usage     string
	Required  bool
	Pattern   string
}

// GetName return name of the ctl flag
func (cf BaseCtlFlag) GetName() string {
	return cf.Name
}

// GetRequired return whether the ctl flag is required
func (cf BaseCtlFlag) GetRequired() bool {
	return cf.Required
}

// StringCtlFlag used for string flag infos
type StringCtlFlag struct {
	BaseCtlFlag
	Value string
}

// NewStringCtlFlag return a StringCtlFlag with given attrs
func NewStringCtlFlag(name, shortHand, usage, value string) StringCtlFlag {
	f := StringCtlFlag{}
	f.Name, f.ShortHand, f.Usage, f.Value = name, shortHand, usage, value
	return f
}

// SetRequired set a StringCtlFlag as required
func (cf StringCtlFlag) SetRequired() FlagImp {
	cf.Required = true
	cf.Usage = cf.Usage + " (required)"
	return cf
}

// SetPattern set a StringCtlFlag's pattern for regexp check
func (cf StringCtlFlag) SetPattern(p string) FlagImp {
	cf.Pattern = p
	return cf
}

// CheckRequired checks the value is not the default type value, var ctl flag checks not blank string
func (cf StringCtlFlag) CheckRequired(v string) bool {
	return v != ""
}

// StringVarP return StringCtlFlag's attribute as StringVarP format
func (cf StringCtlFlag) StringVarP(p *string) (*string, string, string, string, string) {
	return p, cf.Name, cf.ShortHand, cf.Value, cf.Usage
}

// StringVar return StringCtlFlag's attribute as StringVar format
func (cf StringCtlFlag) StringVar(p *string) (*string, string, string, string) {
	return p, cf.Name, cf.Value, cf.Usage
}

// BoolCtlFlag used for bool flag infos
type BoolCtlFlag struct {
	BaseCtlFlag
	Value bool
}

// NewBoolCtlFlag return a BoolCtlFlag with given attrs
func NewBoolCtlFlag(name, shortHand, usage string, value bool) BoolCtlFlag {
	f := BoolCtlFlag{}
	f.Name, f.ShortHand, f.Usage, f.Value = name, shortHand, usage, value
	return f
}

// SetRequired set a BoolCtlFlag as required
func (cf BoolCtlFlag) SetRequired() FlagImp {
	cf.Required = true
	cf.Usage = cf.Usage + " (required)"
	return cf
}

// SetPattern set a StringCtlFlag's pattern for regexp check
func (cf BoolCtlFlag) SetPattern(p string) FlagImp {
	cf.Pattern = p
	return cf
}

// CheckRequired checks the value is not the default type value, var ctl flag checks not blank string
func (cf BoolCtlFlag) CheckRequired(v string) bool {
	return v != "false"
}

// BoolVarP return BoolCtlFlag's attribute as BoolVarP format
func (cf BoolCtlFlag) BoolVarP(p *bool) (*bool, string, string, bool, string) {
	return p, cf.Name, cf.ShortHand, cf.Value, cf.Usage
}
