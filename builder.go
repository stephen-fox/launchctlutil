package launchctlutil

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	oneIndent    = "    "
	twoIndents   = oneIndent + oneIndent
	threeIndents = twoIndents + oneIndent

	openPlist   = "<plist version=\"1.0\">\n"
	closePlist  = "</plist>\n"
	openKey     = "<key>"
	closeKey    = "</key>\n"
	openDict    = "<dict>\n"
	closeDict   = "</dict>\n"
	openArray   = "<array>\n"
	closeArray  = "</array>\n"
	openString  = "<string>"
	closeString = "</string>\n"
	openInt     = "<integer>"
	closeInt    = "</integer>\n"
	newLine     = "\n"
)

// ConfigurationBuilder is used to build a new launchd service Configuration.
//
// Example:
//	config, err := launchctlutil.NewConfigurationBuilder().
//		SetKind(launchctlutil.UserAgent).
//		SetLabel("com.testing").
//		SetRunAtLoad(true).
//		SetCommand("echo").
//		AddArgument("Hello world!").
//		SetLogParentPath("/tmp").
//		Build()
//	if err != nil {
//		log.Fatal(err.Error())
//	}
type ConfigurationBuilder interface {
	// SetLabel sets the label.
	SetLabel(label string) ConfigurationBuilder

	// SetCommand sets the command to execute.
	SetCommand(command string) ConfigurationBuilder

	// AddEnvironmentVariable adds an environment variable with
	// a given value.
	AddEnvironmentVariable(name string, value string) ConfigurationBuilder

	// AddArgument adds an argument for a command.
	AddArgument(value string) ConfigurationBuilder

	// SetLogParentPath sets the directory path where a log
	// file will be saved to. One combined log file is saved
	// containing the output of both stderr and stdout. The
	// file name is formatted as "(launchd-label).log".
	//
	// This setting overrides the settings of SetStandardErrorPath()
	// and SetStandardOutPath().
	SetLogParentPath(logParentPath string) ConfigurationBuilder

	// SetStandardErrorPath sets the file path where stderr
	// output will be saved to.
	//
	// This setting is ignored if SetLogParentPath() is used.
	SetStandardErrorPath(filePath string) ConfigurationBuilder

	// SetStandardOutPath sets the file path where stdout
	// output will be saved to.
	//
	// This setting is ignored if SetLogParentPath() is used.
	SetStandardOutPath(filePath string) ConfigurationBuilder

	// SetKind sets the type.
	SetKind(kind Kind) ConfigurationBuilder

	// SetStartInterval sets the start interval in seconds.
	SetStartInterval(seconds int) ConfigurationBuilder

	// SetStartCalendarIntervalMinute sets the minute of each hour
	// that the command will be executed. For example, setting the
	// minute to 10 will run the command at the 10th minute of each
	// hour: 01:10, 02:10, 03:10, and so on.
	SetStartCalendarIntervalMinute(minuteOfEachHour int) ConfigurationBuilder

	// SetRunAtLoad sets whether or not the service will start
	// when it is loaded.
	SetRunAtLoad(enabled bool) ConfigurationBuilder

	// SetUserName sets whether the service should run as a specific
	// user (by username).
	SetUserName(userName string) ConfigurationBuilder

	// SetGroupName sets whether the service should run as a specific
	// group (by group name).
	SetGroupName(groupName string) ConfigurationBuilder

	// SetInitGroups sets whether launchd should call the function
	// initgroups(3) before starting the servie.
	SetInitGroups(enabled bool) ConfigurationBuilder

	// SetUmask sets the umask for the service.
	SetUmask(umask int) ConfigurationBuilder

	// Build returns the resulting service Configuration.
	Build() (Configuration, error)
}

type configurationBuilder struct {
	label                             string
	command                           string
	environmentVariables              string
	arguments                         string
	lines                             string
	logParentPath                     string
	stderrLogFilePath                 string
	stdoutLogFilePath                 string
	configurationFilePath             string
	kind                              Kind
	startIntervalSeconds              int
	startCalendarIntervalMinuteOfHour int
	isStartCalendarIntervalMinuteSet  bool
	runAtLoad                         bool
	isRunAtLoadSet                    bool
	userName                          string
	groupName                         string
	initGroups                        bool
	isInitGroupsSet                   bool
	umask                             int
	isUmaskSet                        bool
}

// NewConfigurationBuilder creates a new instance of a ConfigurationBuilder.
func NewConfigurationBuilder() ConfigurationBuilder {
	return &configurationBuilder{}
}

func (o *configurationBuilder) SetLabel(label string) ConfigurationBuilder {
	o.label = label
	return o
}

func (o *configurationBuilder) SetCommand(command string) ConfigurationBuilder {
	o.command = command
	return o
}

func (o *configurationBuilder) AddEnvironmentVariable(name string, value string) ConfigurationBuilder {
	o.environmentVariables = concat(o.environmentVariables, "            ", openKey, name, closeKey,
		threeIndents, openString, value, closeString)
	return o
}

func (o *configurationBuilder) AddArgument(value string) ConfigurationBuilder {
	o.arguments = concat(o.arguments,
		threeIndents, openString, value, closeString)
	return o
}

func (o *configurationBuilder) SetLogParentPath(logParentPath string) ConfigurationBuilder {
	o.logParentPath = logParentPath
	return o
}

func (o *configurationBuilder) SetStandardErrorPath(filePath string) ConfigurationBuilder {
	o.stderrLogFilePath = filePath
	return o
}

func (o *configurationBuilder) SetStandardOutPath(filePath string) ConfigurationBuilder {
	o.stdoutLogFilePath = filePath
	return o
}

func (o *configurationBuilder) SetKind(kind Kind) ConfigurationBuilder {
	o.kind = kind
	return o
}

func (o *configurationBuilder) SetStartInterval(seconds int) ConfigurationBuilder {
	o.startIntervalSeconds = seconds
	return o
}

func (o *configurationBuilder) SetStartCalendarIntervalMinute(minuteOfEachHour int) ConfigurationBuilder {
	o.startCalendarIntervalMinuteOfHour = minuteOfEachHour
	o.isStartCalendarIntervalMinuteSet = true
	return o
}

func (o *configurationBuilder) SetRunAtLoad(enabled bool) ConfigurationBuilder {
	o.runAtLoad = enabled
	o.isRunAtLoadSet = true
	return o
}

func (o *configurationBuilder) SetUserName(userName string) ConfigurationBuilder {
	o.userName = userName
	return o
}

func (o *configurationBuilder) SetGroupName(groupName string) ConfigurationBuilder {
	o.groupName = groupName
	return o
}

func (o *configurationBuilder) SetInitGroups(enabled bool) ConfigurationBuilder {
	o.initGroups = enabled
	o.isInitGroupsSet = true
	return o
}

func (o *configurationBuilder) SetUmask(umask int) ConfigurationBuilder {
	o.umask = umask
	o.isUmaskSet = true
	return o
}

func (o *configurationBuilder) Build() (Configuration, error) {
	lines := concat("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n",
		"<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" ",
		"\"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n",
		openPlist,
		oneIndent, openDict)

	lines = concat(lines, twoIndents, openKey, "Label", closeKey,
		twoIndents, openString, o.label, closeString)

	if len(o.environmentVariables) > 0 {
		lines = concat(lines, twoIndents, openKey, "EnvironmentVariables", closeKey,
			twoIndents, openDict,
			o.environmentVariables,
			twoIndents, closeDict)
	}

	if len(o.userName) > 0 {
		lines = concat(lines, twoIndents, openKey, "UserName", closeKey,
			twoIndents, openString, o.userName, closeString)
	}

	if len(o.groupName) > 0 {
		lines = concat(lines, twoIndents, openKey, "GroupName", closeKey,
			twoIndents, openString, o.groupName, closeString)
	}

	if o.isInitGroupsSet {
		lines = concat(lines, twoIndents, openKey, "InitGroups", closeKey,
			twoIndents, boolToXml(o.initGroups), newLine)
	}

	if o.isUmaskSet {
		lines = concat(lines, twoIndents, openKey, "Umask", closeKey,
			twoIndents, openInt, strconv.Itoa(o.umask), closeInt)
	}

	if len(o.command) > 0 {
		lines = concat(lines, twoIndents, openKey, "ProgramArguments", closeKey,
			twoIndents, openArray,
			threeIndents, openString, o.command, closeString)

		if len(o.arguments) > 0 {
			lines = concat(lines, o.arguments)
		}

		lines = concat(lines, twoIndents, closeArray)
	}

	if len(o.logParentPath) > 0 {
		lines = concat(lines, twoIndents, openKey, "StandardOutPath", closeKey,
			twoIndents, openString, o.logParentPath, "/", o.label, ".log", closeString)

		lines = concat(lines, twoIndents, openKey, "StandardErrorPath", closeKey,
			twoIndents, openString, o.logParentPath, "/", o.label, ".log", closeString)
	} else {
		if len(o.stderrLogFilePath) > 0 {
			lines = concat(lines, twoIndents, openKey, "StandardErrorPath", closeKey,
				twoIndents, openString, o.stderrLogFilePath, closeString)
		}

		if len(o.stdoutLogFilePath) > 0 {
			lines = concat(lines, twoIndents, openKey, "StandardOutPath", closeKey,
				twoIndents, openString, o.stdoutLogFilePath, closeString)
		}
	}

	if o.startIntervalSeconds > 0 {
		lines = concat(lines, twoIndents, openKey, "StartInterval", closeKey,
			twoIndents, openInt, strconv.Itoa(o.startIntervalSeconds), closeInt)

	}

	if o.isStartCalendarIntervalMinuteSet {
		lines = concat(lines, twoIndents, openKey, "StartCalendarInterval", closeKey,
			twoIndents, openDict,
			threeIndents, openKey, "Minute", closeKey,
			threeIndents, openInt, strconv.Itoa(o.startCalendarIntervalMinuteOfHour), closeInt,
			twoIndents, closeDict)
	}

	if o.isRunAtLoadSet {
		lines = concat(lines, twoIndents, openKey, "RunAtLoad", closeKey,
			twoIndents, boolToXml(o.runAtLoad), newLine)
	}

	lines = concat(lines, oneIndent, closeDict, closePlist)

	return &configuration{
		label:    o.label,
		contents: lines,
		kind:     o.kind,
	}, nil
}

func boolToXml(b bool) string {
	return fmt.Sprintf("<%t/>", b)
}

func concat(current string, additions ...string) (new string) {
	var buffer bytes.Buffer

	buffer.WriteString(current)
	for _, addition := range additions {
		buffer.WriteString(addition)
	}

	return buffer.String()
}
