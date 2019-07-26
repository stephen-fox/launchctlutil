package launchctlutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

const (
	Unknown      Status = "unknown"
	NotInstalled Status = "not_installed"
	Running      Status = "running"
	NotRunning   Status = "not_running"
)

// Status represents the status of a launchd service.
type Status string

// StatusDetails provides detailed information about the status of
// a launchd service.
type StatusDetails struct {
	Status            Status
	Pid               int
	LastExitStatus    int
	PidErr            error
	LastExitStatusErr error
}

// GotLastExitStatus returns true if the launchd service provided an
// exit status.
func (o StatusDetails) GotLastExitStatus() bool {
	return o.LastExitStatusErr == nil
}

// GotPid returns true if the launchd service provided a PID.
func (o StatusDetails) GotPid() bool {
	return o.PidErr == nil
}

const (
	defaultLaunchctl          = "launchctl"
	couldNotFindServicePrefix = "Could not find service "
	lastExitStatusPrefix      = "\"LastExitStatus\" = "
	pidPrefix                 = "\"PID\" = "
	serviceListLineSuffix     = ";"
)

var (
	// ExePath is the path to the launchctl CLI application.
	ExePath = defaultLaunchctl
)

// Install installs the provided service Configuration.
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
		return fmt.Errorf("an unknown error occurred installing the laucnctl config")
	}

	return nil
}

// Remove unloads and removes the specified service configuration file.
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

// IsInstalled is a wrapper for Configuration.IsInstalled().
func IsInstalled(configuration Configuration) (isInstalled bool, err error) {
	return configuration.IsInstalled()
}

// Start starts the specified launchd service.
func Start(label string, kind Kind) error {
	if kind == Daemon {
		err := isRoot()
		if err != nil {
			return err
		}
	}

	_, err := run("start", label)
	if err != nil {
		return err
	}

	return nil
}

// Stop stops the specified launchd service.
func Stop(label string, kind Kind) error {
	if kind == Daemon {
		err := isRoot()
		if err != nil {
			return err
		}
	}

	_, err := run("stop", label)
	if err != nil {
		return err
	}

	return nil
}

// CurrentStatus returns the current status of the specified launchd service.
func CurrentStatus(label string) (StatusDetails, error) {
	output, err := run("list", label)
	if err != nil {
		if strings.HasPrefix(output, couldNotFindServicePrefix) {
			return StatusDetails{
				Status: NotInstalled,
			}, nil
		}

		return StatusDetails{
			Status: Unknown,
		}, err
	}

	details := StatusDetails{
		Status: NotRunning,
	}

	for _, l := range strings.Split(output, "\n") {
		l = strings.TrimSpace(l)

		if strings.HasPrefix(l, lastExitStatusPrefix) {
			exit, err := getLastExitStatus(l)
			if err != nil {
				details.LastExitStatusErr = err
				continue
			}

			details.LastExitStatus = exit
		}

		if strings.HasPrefix(l, pidPrefix) {
			pid, err := getPid(l)
			if err != nil {
				details.PidErr = err
				continue
			}

			details.Pid = pid
			details.Status = Running
		}
	}

	return details, nil
}

func getPid(lineWithoutLeadingSpaces string) (int, error) {
	lineWithoutLeadingSpaces = strings.TrimPrefix(lineWithoutLeadingSpaces, pidPrefix)
	lineWithoutLeadingSpaces = strings.TrimSuffix(lineWithoutLeadingSpaces, serviceListLineSuffix)

	pid, err := strconv.Atoi(lineWithoutLeadingSpaces)
	if err != nil {
		return 0, err
	}

	return pid, nil
}

func getLastExitStatus(lineWithoutLeadingSpaces string) (int, error) {
	lineWithoutLeadingSpaces = strings.TrimPrefix(lineWithoutLeadingSpaces, lastExitStatusPrefix)
	lineWithoutLeadingSpaces = strings.TrimSuffix(lineWithoutLeadingSpaces, serviceListLineSuffix)

	exit, err := strconv.Atoi(lineWithoutLeadingSpaces)
	if err != nil {
		return 0, err
	}

	return exit, nil
}

func isRoot() error {
	currentUser, err := user.Current()
	if err != nil {
		// For whatever reason, 'user.Current()' throws a "not
		// implemented error" when running as a launch service
		// on macOS.
		if runtime.GOOS == "darwin" && strings.Contains(err.Error(), "Current not implemented on") {
			return nil
		}
		return fmt.Errorf("failed to check if current user is root - %s", err.Error())
	}

	if currentUser.Username == "root" {
		return nil
	}

	return fmt.Errorf("root privileges are required to do this")
}

func run(args... string) (output string, err error) {
	command := exec.Command(ExePath, args...)
	raw, err := command.CombinedOutput()
	output = string(raw)
	if err != nil {
		return output, fmt.Errorf("%s - output: %s", err.Error(), output)
	}

	if strings.Contains(output, ": Invalid property list") {
		return output, fmt.Errorf("invalid property list")
	}

	return output, nil
}
