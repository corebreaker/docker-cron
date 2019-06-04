package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/gorhill/cronexpr"
	gerr "github.com/corebreaker/goerrors"
)

var (
	command string
	args []string
)

func init() {
	choice := os.Getenv("COMMAND_TYPE")
	if choice == "" {
		choice = "composer"
	}

	switch choice {
	case "standard":
		command = "docker"
		args = []string{"exec", "-ti"}

	case "composer":
		command = "docker-compose"
		args = []string{"exec"}

	default:
		LogFatal("%s", gerr.MakeError("Bad command type; valid choices are: `standard`, `composer`"))
	}
}

type Job struct {
	At        string               `json:"at" yaml:"at"`
	Container string               `json:"container" yaml:"container"`
	Command   string               `json:"command" yaml:"command"`
	expr      *cronexpr.Expression `json:"-" yaml:"-"`
	next      time.Time            `json:"-" yaml:"-"`
}

func (j Job) Run() {
	cmd := exec.Command(command, append(args, j.Container, j.Command)...)

	logerr := func(err error) {
		LogError("Command error on %s (%s): %s", j.Container, j.Command, err)
		LogDebug("Command: %s %s", cmd.Path, cmd.Args)
		LogDebug("Process State: %s", cmd.ProcessState)
		LogDebug("CWD:", cmd.Dir)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logerr(gerr.DecorateError(err))

		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logerr(gerr.DecorateError(err))

		return
	}

	if err := cmd.Start(); err != nil {
		logerr(err)
	}

	MakeOutputPipe(stdout, "OUT", j.Container).Start()
	MakeOutputPipe(stderr, "ERR", j.Container).Start()

	if err := cmd.Wait(); err != nil {
		logerr(gerr.DecorateError(err))
	}
}
