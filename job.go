package main

import (
	"os"
	"os/exec"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/corebreaker/goerrors"
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
		LogFatal("%s", goerrors.MakeError("Bad command type; valid choices are: `standard`, `composer`"))
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
	logerr := func(err error) {
		LogError("Command error on %s (%s): %s", j.Container, j.Command, err)
	}

	cmd := exec.Command(command, append(args, j.Container, j.Command)...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logerr(err)

		return
	}

	stderr, err := cmd.StdoutPipe()
	if err != nil {
		logerr(err)

		return
	}

	MakeOutputPipe(stdout, cmd, "OUT", j.Container).Start()
	MakeOutputPipe(stderr, cmd, "ERR", j.Container).Start()

	if err := cmd.Run(); err != nil {
		logerr(err)
	}
}
