package osquery

import (
	"github.com/stretchr/testify/mock"
)

type MockOsqueryClient struct {
	mock.Mock
}

func (m *MockOsqueryClient) ExecuteQuery(query string) ([]map[string]string, error) {
	args := m.Called(query)
	return args.Get(0).([]map[string]string), args.Error(1)
}

func (m *MockOsqueryClient) GetOsVersion() (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}

func (m *MockOsqueryClient) GetOsqueryVersion() (string, error) {
	args := m.Called()
	return args.Get(0).(string), args.Error(1)
}

func (m *MockOsqueryClient) GetAppsInstalled() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}
