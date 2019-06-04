package main

import (
	"os"
	"runtime"
	"strings"

	gerr "github.com/corebreaker/goerrors"
)

func main() {
	switch strings.ToLower(os.Getenv("DEBUG")) {
	case "1", "on", "yes", "y", "o":
		gerr.SetDebug(true)
	}

	StartScheduler()

	if err := InitAutoUpdate(); err != nil {
		LogError("%s", err)

		return
	}

	runtime.Goexit()
}
