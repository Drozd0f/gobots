package ytdl

import (
	"io"
	"os/exec"
)

type stream struct {
	cmd *exec.Cmd
}

func (p *stream) GetOutput() (io.ReadCloser, error) {
	return p.cmd.StdoutPipe()
}

func (p *stream) Start() error {
	return p.cmd.Start()
}

func (p *stream) Cancel() error {
	return p.cmd.Process.Kill()
}

func (p *stream) Wait() error {
	return p.cmd.Wait()
}
