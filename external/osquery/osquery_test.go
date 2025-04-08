package osquery

import (
	"errors"
	"github.com/sourava/secfix/external/command_executor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestExecuteQuery_ShouldReturnError_WhenExecCommandReturnsError(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte(""), errors.New("some error"))

	client := NewOsqueryClient(mockExecutor)
	res, err := client.ExecuteQuery("")

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestExecuteQuery_ShouldReturnError_WhenExecCommandReturnsInvalidResponse(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("{}"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.ExecuteQuery("")

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestExecuteQuery_ShouldReturnError_WhenExecCommandReturns0Rows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.ExecuteQuery("")

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestExecuteQuery_ShouldReturnResponse_WhenExecCommandReturnsMoreThan0Rows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[{\"version\":\"1.0.0\"}]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.ExecuteQuery("")

	assert.Equal(t, []map[string]string{{"version": "1.0.0"}}, res)
	assert.Nil(t, err)
}

func TestGetOsVersion_ShouldReturnError_WhenVersionNotFoundInRows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[{\"abc\":\"1.0.0\"}]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.GetOsVersion()

	assert.Equal(t, "", res)
	assert.NotNil(t, err)
}

func TestGetOsVersion_ShouldReturnVersion_WhenVersionFoundInRows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[{\"version\":\"1.0.0\"}]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.GetOsVersion()

	assert.Equal(t, "1.0.0", res)
	assert.Nil(t, err)
}

func TestGetOsqueryVersion_ShouldReturnError_WhenVersionNotFoundInRows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[{\"abc\":\"1.0.0\"}]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.GetOsqueryVersion()

	assert.Equal(t, "", res)
	assert.NotNil(t, err)
}

func TestGetOsqueryVersion_ShouldReturnVersion_WhenVersionFoundInRows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[{\"version\":\"1.0.0\"}]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.GetOsqueryVersion()

	assert.Equal(t, "1.0.0", res)
	assert.Nil(t, err)
}

func TestGetAppsInstalled_ShouldReturnAppsInstalled_WhenNameFoundInRows(t *testing.T) {
	mockExecutor := new(command_executor.MockCommandExecutor)
	mockExecutor.On("CombinedOutput", mock.Anything).Return([]byte("[{\"name\":\"App1\"},{\"name\":\"App2\"}]"), nil)

	client := NewOsqueryClient(mockExecutor)
	res, err := client.GetAppsInstalled()

	assert.Equal(t, []string{"App1", "App2"}, res)
	assert.Nil(t, err)
}
