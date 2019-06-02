package main

import (
	"github.com/corebreaker/goerrors"
	"os"
	"os/exec"
	"time"

	"github.com/gorhill/cronexpr"
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
	cmd := exec.Command(command, append(args, j.Container, j.Command)...)
	if err := cmd.Run(); err != nil {
		LogError("Command error on %s (%s): %s", j.Container, j.Command, err)
	}
}
