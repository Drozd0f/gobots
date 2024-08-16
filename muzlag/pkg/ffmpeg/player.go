package ffmpeg

import (
	"bufio"
	"io"
	"os/exec"
)

type player struct {
	cmd        *exec.Cmd
	bufferSize int
}

func (p *player) SetInput(r io.Reader) {
	p.cmd.Stdin = bufio.NewReaderSize(r, p.bufferSize)
}

func (p *player) GetOutput() (io.ReadCloser, error) {
	return p.cmd.StdoutPipe()
}

func (p *player) Start() error {
	return p.cmd.Start()
}

func (p *player) Cancel() error {
	return p.cmd.Process.Kill()
}

func (p *player) Wait() error {
	return p.cmd.Wait()
}
