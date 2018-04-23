package launchctlutil

import (
	"errors"
	"os"
)

const (
	UserAgent Kind = iota
	Daemon    Kind = iota
)

type Kind int

type Configuration interface {
	GetLabel() string

	GetContents() string

	GetFilePath() (configFilePath string, err error)

	GetKind() Kind
}

type configuration struct {
	label    string
	contents string
	kind     Kind
}

func (c *configuration) GetLabel() string {
	return c.label
}

func (c *configuration) GetContents() string {
	return c.contents
}

func (c *configuration) GetFilePath() (configFilePath string, err error) {
	configFilePath = ""

	switch c.kind {
	case UserAgent:
		homePath := os.Getenv("HOME")
		if homePath == "" {
			return "", errors.New("Failed to determine HOME for UserAgent launchctl configuration")
		}
		configFilePath = homePath + "/Library/LaunchAgents"
	case Daemon:
		configFilePath = "/Library/LaunchDaemons"
	default:
		return "", errors.New("An unknown launchctl configuration type was specified")
	}

	return configFilePath + "/" + c.label + ".plist", nil
}

func (c *configuration) GetKind() Kind {
	return c.kind
}