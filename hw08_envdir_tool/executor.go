package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func runCmd(writer io.Writer, cmd []string, env Environment) int {
	ind := 0
	environ := make([]string, len(env))
	for key, item := range env {
		environ[ind] = fmt.Sprintf("%s=%s", key, item.Value)
		ind++
	}

	command := exec.Cmd{
		Path:   cmd[0],
		Args:   cmd[0:],
		Env:    environ,
		Stdout: writer,
	}

	if err := command.Run(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			return exit.ExitCode()
		}
	}

	return 0
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	return runCmd(os.Stdout, cmd, env)
}
