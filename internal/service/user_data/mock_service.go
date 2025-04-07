package user_data

import (
	"github.com/sourava/secfix/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserDataService struct {
	mock.Mock
}

func (m *MockUserDataService) GetLatestUserData() (*models.UserData, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserData), args.Error(1)
}

func (m *MockUserDataService) AddLatestUserData(data *models.UserData) error {
	args := m.Called(data)
	return args.Error(0)
}
