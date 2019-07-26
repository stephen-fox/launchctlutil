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
		t.Fatalf("status should be running - got %s", details.Status)
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
		t.Fatalf("status should be not running - got %s", details.Status)
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
		t.Fatalf("exit status should be %d - got %d", exp, exit)
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
		t.Fatalf("PID should be %d - got %d", exp, pid)
	}
}

func TestCurrentStatusNotInstalled(t *testing.T) {
	details, err := CurrentStatus("com.apple.nevergoingtoexisttrash")
	if err != nil {
		t.Fatal(err.Error())
	}

	if details.Status != NotInstalled {
		t.Fatalf("status should be running - got %s", details.Status)
	}
}

func TestStatusDetails_GotLastExitStatus(t *testing.T) {
	details := StatusDetails{
		LastExitStatusErr: errors.New("broke"),
	}

	if details.GotLastExitStatus() {
		t.Fatal("got last exit status when there was an error")
	}
}

func TestStatusDetails_GotPid(t *testing.T) {
	details := StatusDetails{
		PidErr: errors.New("broke"),
	}

	if details.GotPid() {
		t.Fatal("got a PID when there was an error")
	}
}
