package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"

	gerr "github.com/corebreaker/goerrors"
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
		log.Printf("%s", gerr.DecorateError(err))
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
		log.Printf("%s", gerr.MakeError("Bad command type; valid choices are: `standard`, `composer`"))
	}
}

type OutputPipe struct {
	in io.ReadCloser
	outType string
	container string
}

func (p OutputPipe) Start() {
	in := bufio.NewReader(p.in)

	tmpl := "[%s] {%s}: %s"

	go func() {
		defer func() { _ = p.in.Close() }()

		for {
			out, err := in.ReadString('\n')
			if err != nil {
				if (err != io.EOF) && (err != os.ErrClosed) {
					pathErr, ok := err.(*os.PathError)

					if (!ok) || ((pathErr.Err != io.EOF) && (pathErr.Err != os.ErrClosed)) {
						log.Printf(tmpl, p.outType, p.container, err)
					}
				}

				return
			}

			log.Printf(tmpl, p.outType, p.container, out)
		}
	}()
}

func MakeOutputPipe(in io.ReadCloser, outType, container string) OutputPipe {
	return OutputPipe{
		in: in,
		outType: outType,
		container: container,
	}
}

func main() {
	nullFile, err := os.Open(os.DevNull)
	if err != nil {
		log.Printf("%s", gerr.DecorateError(err))

		return
	}

	defer func() {
		_ = nullFile.Close()
	}()

	cmd := exec.Command(command, append(args, "web", "bin/cron-unconfirmed-users-purge.sh")...)
	cmd.Stdin = nullFile

	logerr := func(err error) {
		log.Printf("Command error: %s", err)
		log.Printf("Command: %s %s", cmd.Path, cmd.Args)
		log.Printf("Process State: %s", cmd.ProcessState)
		log.Printf("CWD: %s", cmd.Dir)
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

	cmd.Dir = cwd
	if err := cmd.Start(); err != nil {
		logerr(err)
	}

	MakeOutputPipe(stdout, "OUT", "web").Start()
	MakeOutputPipe(stderr, "ERR", "web").Start()

	if err := cmd.Wait(); err != nil {
		logerr(gerr.DecorateError(err))
	}
}
