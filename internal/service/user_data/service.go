package user_data

import (
	"encoding/json"
	"github.com/sourava/secfix/internal/models"
	"gorm.io/gorm"
)

type UserDataServiceInterface interface {
	GetLatestUserData() (*models.UserData, error)
	AddLatestUserData(data *models.UserData) error
}

type UserDataService struct {
	db *gorm.DB
}

func NewUserDataService(db *gorm.DB) *UserDataService {
	return &UserDataService{
		db: db,
	}
}

func (service *UserDataService) GetLatestUserData() (*models.UserData, error) {
	var versionInfo models.VersionInfo
	var installedApplications models.InstalledApplications
	if err := service.db.Order("created_at desc").First(&versionInfo).Error; err != nil {
		return nil, err
	}
	if err := service.db.Order("created_at desc").First(&installedApplications).Error; err != nil {
		return nil, err
	}

	var appsInstalled []string
	err := json.Unmarshal(installedApplications.AppsInstalled, &appsInstalled)
	if err != nil {
		return nil, err
	}
	return &models.UserData{
		OSVersion:      versionInfo.OSVersion,
		OSQueryVersion: versionInfo.OSQueryVersion,
		AppsInstalled:  appsInstalled,
	}, nil
}

func (service *UserDataService) AddLatestUserData(data *models.UserData) error {
	appsInstalled, err := json.Marshal(data.AppsInstalled)
	if err != nil {
		return err
	}

	if err = service.db.Create(&models.VersionInfo{OSQueryVersion: data.OSQueryVersion, OSVersion: data.OSVersion}).Error; err != nil {
		return err
	}
	if err = service.db.Create(&models.InstalledApplications{AppsInstalled: appsInstalled}).Error; err != nil {
		return err
	}
	return nil
}
