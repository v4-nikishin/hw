package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	Seccess = 0
	Err     = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	returnCode = Seccess

	if len(cmd) == 0 {
		fmt.Println("Command is empty")
		return Err
	}

	proc := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	for k, v := range env {
		if v.Value == "" && v.NeedRemove {
			os.Unsetenv(k)
			continue
		}
		if v.NeedRemove {
			os.Unsetenv(k)
		}
		os.Setenv(k, v.Value)
	}
	proc.Env = os.Environ()
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	if err := proc.Run(); err != nil {
		fmt.Println(err)
		returnCode = Err
	}
	return
}
