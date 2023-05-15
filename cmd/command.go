package cmd

import (
	"os/exec"

	"github.com/mattn/go-shellwords"
)

func Run(cmd string) error {
	p := shellwords.NewParser()
	args, err := p.Parse(cmd)
	if err != nil {
		panic(err)
	}
	command := args[0]
	args = args[1:]
	r := exec.Command(command, args...)
	return r.Start()
}
