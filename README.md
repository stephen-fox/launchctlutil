# launchctlutil

## What is it?
A Go library for working with launchctl.

## API
The library offers a basic API for interacting with launchd via launchctl.
This includes building configurations and installing them:
```go
package main

import (
	"fmt"
	"log"

	"github.com/stephen-fox/launchctlutil"
)

func main() {
	config, err := launchctlutil.NewConfigurationBuilder().
		SetKind(launchctlutil.UserAgent).
		SetLabel("com.testing").
		SetRunAtLoad(true).
		SetCommand("echo").
		AddArgument("Hello world!").
		SetLogParentPath("/tmp").
		Build()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Configuration contents:\n" + config.GetContents())

	err = launchctlutil.Install(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Remove by running:
	// launchctl remove com.testing
	// rm ~/Library/LaunchAgents/com.testing.plist 
}
```

You can also get the status of an agent (or a daemon if running as `root`):
```go
package main

import (
	"fmt"
	"log"

	"github.com/stephen-fox/launchctlutil"
)

func main() {
	details, err := launchctlutil.CurrentStatus("com.apple.Dock.agent")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Status:", details.Status)

	if details.GotPid() {
		log.Println("Current PID:", details.Pid)
	}

	if details.GotLastExitStatus() {
		log.Println("Last exit status:", details.LastExitStatus)
	}
}
```
