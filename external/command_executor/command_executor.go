package command_executor

import "os/exec"

type CommandExecutor interface {
	CombinedOutput(cmd *exec.Cmd) ([]byte, error)
}

type DefaultExecutor struct{}

func (e *DefaultExecutor) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	return cmd.CombinedOutput()
}
