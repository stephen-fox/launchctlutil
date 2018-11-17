package launchctlutil

import (
	"bytes"
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
)

type ConfigurationBuilder interface {
	SetLabel(label string) ConfigurationBuilder

	SetCommand(command string) ConfigurationBuilder

	AddEnvironmentVariable(name string, value string) ConfigurationBuilder

	AddArgument(value string) ConfigurationBuilder

	SetLogParentPath(logParentPath string) ConfigurationBuilder

	SetKind(kind Kind) ConfigurationBuilder

	SetStartInterval(seconds int) ConfigurationBuilder

	SetStartCalendarIntervalMinute(minuteOfEachHour int) ConfigurationBuilder

	SetRunAtLoad(enabled bool) ConfigurationBuilder

	Build() (Configuration, error)
}

type configurationBuilder struct {
	label                             string
	command                           string
	environmentVariables              string
	arguments                         string
	lines                             string
	logParentPath                     string
	configurationFilePath             string
	kind                              Kind
	startIntervalSeconds              int
	startCalendarIntervalMinuteOfHour int
	isStartCalendarIntervalMinuteSet  bool
	runAtLoad                         bool
}

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

	if o.environmentVariables != "" {
		lines = concat(lines, twoIndents, openKey, "EnvironmentVariables", closeKey,
			twoIndents, openDict,
			o.environmentVariables,
			twoIndents, closeDict)
	}

	if o.command != "" {
		lines = concat(lines, twoIndents, openKey, "ProgramArguments", closeKey,
			twoIndents, openArray,
			threeIndents, openString, o.command, closeString)

		if o.arguments != "" {
			lines = concat(lines, o.arguments)
		}

		lines = concat(lines, twoIndents, closeArray)
	}

	if o.logParentPath != "" {
		lines = concat(lines, twoIndents, openKey, "StandardOutPath", closeKey,
			twoIndents, openString, o.logParentPath, "/", o.label, ".log", closeString,)

		lines = concat(lines, twoIndents, openKey, "StandardErrorPath", closeKey,
			twoIndents, openString, o.logParentPath, "/", o.label, ".log", closeString)
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

	if o.runAtLoad {
		lines = concat(lines, twoIndents, openKey, "RunAtLoad", closeKey,
			twoIndents, "<true/>\n")
	}

	lines = concat(lines, oneIndent, closeDict, closePlist)

	return &configuration{
		label:    o.label,
		contents: lines,
		kind:     o.kind,
	}, nil
}

func concat(current string, additions ...string) (new string) {
	var buffer bytes.Buffer

	buffer.WriteString(current)
	for _, addition := range additions {
		buffer.WriteString(addition)
	}

	return buffer.String()
}
