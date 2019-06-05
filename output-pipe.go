package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type OutputPipe struct {
	in        io.ReadCloser
	outType   string
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
						LogError(tmpl, p.outType, p.container, err)
					}
				}

				return
			}

			LogInfo(tmpl, p.outType, p.container, strings.TrimRight(out, " \r\n"))
		}
	}()
}

func MakeOutputPipe(in io.ReadCloser, outType, container string) OutputPipe {
	return OutputPipe{
		in:        in,
		outType:   outType,
		container: container,
	}
}
