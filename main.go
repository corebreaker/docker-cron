package main

import (
	"os"
	"runtime"
	"strings"

	gerr "github.com/corebreaker/goerrors"
)

func mainHandler() error {
	switch strings.ToLower(os.Getenv("DEBUG")) {
	case "1", "on", "yes", "y", "o":
		gerr.SetDebug(true)
	}

	StartScheduler()

	if err := InitAutoUpdate(); err != nil {
		LogError("%s", err)

		return err
	}

	runtime.Goexit()

	return nil
}

func main() {
	gerr.CheckedMain(mainHandler)
}
