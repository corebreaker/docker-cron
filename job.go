package main

import (
	"os"
	"os/exec"
	"time"

	gerr "github.com/corebreaker/goerrors"
	"github.com/gorhill/cronexpr"
)

var (
	cwd string
	command string
	args []string
)

func init() {
	var err error

	cwd, err = os.Getwd()
	if err != nil {
		LogFatal("%s", gerr.DecorateError(err))
	}

	choice := os.Getenv("COMMAND_TYPE")
	if choice == "" {
		choice = "composer"
	}

	switch choice {
	case "standard":
		command = "docker"
		args = []string{"exec", "-t"}

	case "composer":
		command = "docker-compose"
		args = []string{"exec", "-T"}

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

	nullFile, err := os.Open(os.DevNull)
	if err != nil {
		LogError("Command error on %s (%s): %s", j.Container, j.Command, err)
		LogDebug("Command: %s %s", cmd.Path, cmd.Args)

		return
	}

	defer func() { _ = nullFile.Close() }()

	cmd.Stdin = nullFile
	cmd.Dir = cwd

	logerr := func(err error) {
		LogError("Command error on %s (%s): %s", j.Container, j.Command, err)
		LogDebug("Command: %s %s", cmd.Path, cmd.Args)
		LogDebug("Process State: %s", cmd.ProcessState)
		LogDebug("CWD: %s", cmd.Dir)
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
