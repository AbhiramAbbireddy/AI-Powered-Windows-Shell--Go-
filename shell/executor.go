package shell

import (
	"bytes"
	"os/exec"
)

type ExecutionResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func ExecuteCommand(shellName, command string) (ExecutionResult, error) {
	cmd := exec.Command(shellName, "/C", command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := ExecutionResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if cmd.ProcessState != nil {
		result.ExitCode = cmd.ProcessState.ExitCode()
	}

	return result, err
}
