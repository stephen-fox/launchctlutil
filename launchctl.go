package launchctlutil

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

const (
	defaultLaunchctl = "launchctl"
)

var (
	ExePath = defaultLaunchctl
)

func Install(configuration Configuration) error {
	if configuration.GetKind() == Daemon {
		err := isRoot()
		if err != nil {
			return err
		}
	}

	configPath, err := configuration.GetFilePath()
	if err != nil {
		return err
	}

	// Try to remove the LaunchAgent first because it may already exist.
	// Ignore errors because this may create false positives.
	Remove(configPath, configuration.GetKind())

	err = ioutil.WriteFile(configPath, []byte(configuration.GetContents()), 0600)
	if err != nil {
		return err
	}

	_, err = run("load", configPath)
	if err != nil {
		return err
	}

	// Check that the LaunchAgent was installed using special logic because
	// launchctl seems to return exit status 0 even when an error occurs.
	isInstalled, err := IsInstalled(configuration)
	if err != nil {
		return err
	}

	if !isInstalled {
		// Try to remove the config file if the installation fails.
		// Ignore errors because this may create false positives.
		os.Remove(configPath)
		return errors.New("An unknown error occurred installing the laucnctl config")
	}

	return nil
}

func Remove(configPath string, kind Kind) error {
	if kind == Daemon {
		err := isRoot()
		if err != nil {
			return err
		}
	}

	_, err := run("unload", configPath)
	if err != nil {
		return err
	}

	err = os.Remove(configPath)
	if err != nil {
		return err
	}

	return nil
}

func IsInstalled(configuration Configuration) (isInstalled bool, err error) {
	if configuration.GetKind() == Daemon {
		err := isRoot()
		if err != nil {
			return false, err
		}
	}

	output, err := run("list")
	if err != nil {
		return false, err
	}

	if strings.Contains(output, configuration.GetLabel()) {
		configFilePath, err := configuration.GetFilePath()
		if err != nil {
			return false, err
		}
		_, temp := os.Stat(configFilePath)
		if temp == nil {
			currentContents, err := ioutil.ReadFile(configFilePath)
			if err == nil {
				if string(currentContents) == configuration.GetContents() {
					return true, nil
				}
			} else {
				return false, err
			}
		}
	}

	return false, nil
}

func isRoot() error {
	currentUser, err := user.Current()
	if err != nil {
		// For whatever reason, 'user.Current()' throws a "not implemented
		// error" when running as a launch daemon on macOS.
		if runtime.GOOS == "darwin" && strings.Contains(err.Error(), "Current not implemented on") {
			return nil
		}
		return errors.New("Failed to check if current user is root - " + err.Error())
	}

	if currentUser.Username == "root" {
		return nil
	}

	return errors.New("Root privileges are required to do this")
}

func run(args... string) (output string, err error) {
	command := exec.Command(ExePath, args...)
	raw, err := command.CombinedOutput()
	output = string(raw)
	if err != nil {
		return output, errors.New(err.Error() + " - Output: " + output)
	}

	if strings.Contains(output, ": Invalid property list") {
		return output, errors.New("Invalid property list")
	}

	return output, nil
}