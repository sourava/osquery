package command_executor

import (
	"github.com/stretchr/testify/mock"
	"os/exec"
)

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) CombinedOutput(cmd *exec.Cmd) ([]byte, error) {
	args := m.Called(cmd)
	return args.Get(0).([]byte), args.Error(1)
}
