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
}
