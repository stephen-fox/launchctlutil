package launchctlutil

import (
	"errors"
	"strconv"
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

func TestCurrentStatusNotRunning(t *testing.T) {
	details, err := CurrentStatus("com.apple.DiagnosticReportCleanup.plist")
	if err != nil {
		t.Fatal(err.Error())
	}

	if details.Status == Running {
		t.Fatal("Status should be not running. Got -", details.Status)
	}

	if details.Pid > 0 {
		t.Fatal("PID should not be greater than 0")
	}
}

func TestGetLastExitStatus(t *testing.T) {
	exp := 15
	l := `"LastExitStatus" = ` + strconv.Itoa(exp) + ";"

	exit, err := getLastExitStatus(l)
	if err != nil {
		t.Fatal(err.Error())
	}

	if exit != exp {
		t.Fatal("Exit status should be", exp, "- Got", exit)
	}
}

func TestGetPid(t *testing.T) {
	exp := 33385
	l := `"PID" = ` + strconv.Itoa(exp) + ";"

	pid, err := getPid(l)
	if err != nil {
		t.Fatal(err.Error())
	}

	if pid != exp {
		t.Fatal("PID should be", exp, "- Got", pid)
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
