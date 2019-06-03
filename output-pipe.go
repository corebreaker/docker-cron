package main

import (
	"bufio"
	"io"
	"os/exec"
)

type OutputPipe struct {
	in io.ReadCloser
	cmd *exec.Cmd
	outType string
	container string
}

func (p OutputPipe) Start() {
	defer func() { _ = p.in.Close() }()

	in := bufio.NewReader(p.in)
	pid := p.cmd.ProcessState.Pid()

	const tmpl =  "[%s] {%d:%s}: %s"

	go func() {
		for {
			out, err := in.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					LogError(tmpl, p.outType, pid, p.container, err)
				}

				return
			}

			LogInfo(tmpl, p.outType, pid, p.container, out)
		}
	}()
}

func MakeOutputPipe(in io.ReadCloser, cmd *exec.Cmd, outType, container string) OutputPipe {
	return OutputPipe{
		in: in,
		cmd: cmd,
		outType: outType,
		container: container,
	}
}