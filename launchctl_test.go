package launchctlutil

import (
	"errors"
	"testing"
)

func TestCurrentStatus(t *testing.T) {
	details, err := CurrentStatus("com.apple.diskspaced")
	if err != nil {
		t.Fatal(err.Error())
	}

	if details.Status != Running {
		t.Fatal("Status should be running. Got -", details.Status)
	}

	if details.Pid == 0 {
		t.Fatal("PID should be greater than 0")
	}
}

func TestCurrentStatusNotInstalled(t *testing.T) {
	details, err := CurrentStatus("com.apple.nevergoingtoexisttrash")
	if err != nil {
		t.Fatal(err.Error())
	}

	if details.Status != NotInstalled {
		t.Fatal("Status should be running. Got -", details.Status)
	}
}

func TestStatusDetails_GotLastExitStatus(t *testing.T) {
	details := StatusDetails{
		LastExitStatusErr: errors.New("broke"),
	}

	if details.GotLastExitStatus() {
		t.Fatal("Got last exit status when there was an error")
	}
}

func TestStatusDetails_GotPid(t *testing.T) {
	details := StatusDetails{
		PidErr: errors.New("broke"),
	}

	if details.GotPid() {
		t.Fatal("Got a PID when there was an error")
	}
}
