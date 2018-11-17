package launchctlutil

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
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

	IsInstalled() (bool, error)
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

func (c *configuration) IsInstalled() (bool, error) {
	if c.GetKind() == Daemon {
		err := isRoot()
		if err != nil {
			return false, err
		}
	}

	output, err := run("list")
	if err != nil {
		return false, err
	}

	if strings.Contains(output, c.GetLabel()) {
		configFilePath, err := c.GetFilePath()
		if err != nil {
			return false, err
		}
		_, temp := os.Stat(configFilePath)
		if temp == nil {
			currentContents, err := ioutil.ReadFile(configFilePath)
			if err == nil {
				if string(currentContents) == c.GetContents() {
					return true, nil
				}
			} else {
				return false, err
			}
		}
	}

	return false, nil
}
