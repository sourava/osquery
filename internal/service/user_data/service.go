package user_data

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"github.com/sourava/secfix/external/osquery"
	"github.com/sourava/secfix/internal/models"
	"gorm.io/gorm"
)

type UserDataServiceInterface interface {
	GetLatestUserData() (*models.UserData, error)
	AddLatestUserData() error
}

type UserDataService struct {
	db            *gorm.DB
	osQueryClient osquery.OsqueryClientInterface
}

func NewUserDataService(db *gorm.DB, osQueryClient osquery.OsqueryClientInterface) *UserDataService {
	return &UserDataService{
		db:            db,
		osQueryClient: osQueryClient,
	}
}

func (service *UserDataService) GetLatestUserData() (*models.UserData, error) {
	var versionInfo models.VersionInfo
	var installedApplications models.InstalledApplications
	if err := service.db.Order("created_at desc").First(&versionInfo).Error; err != nil {
		log.Error().Err(err).Msg("failed to get latest version info")
		return nil, err
	}
	if err := service.db.Order("created_at desc").First(&installedApplications).Error; err != nil {
		log.Error().Err(err).Msg("failed to get latest installed applications")
		return nil, err
	}

	var appsInstalled []string
	err := json.Unmarshal(installedApplications.AppsInstalled, &appsInstalled)
	if err != nil {
		log.Error().Err(err).Msg("failed to unmarshal installed applications")
		return nil, err
	}
	return &models.UserData{
		OSVersion:      versionInfo.OSVersion,
		OSQueryVersion: versionInfo.OSQueryVersion,
		AppsInstalled:  appsInstalled,
	}, nil
}

func (service *UserDataService) AddLatestUserData() error {
	osVersion, err := service.osQueryClient.GetOsVersion()
	if err != nil {
		log.Error().Err(err).Msg("failed to execute GetOsVersion")
	}

	osqueryVersion, err := service.osQueryClient.GetOsqueryVersion()
	if err != nil {
		log.Error().Err(err).Msg("failed to execute GetOsqueryVersion")
	}

	appsInstalledRes, err := service.osQueryClient.GetAppsInstalled()
	if err != nil {
		log.Error().Err(err).Msg("failed to execute GetAppsInstalled")
	}

	appsInstalled, err := json.Marshal(appsInstalledRes)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal appsInstalled")
		return err
	}

	if err = service.db.Create(&models.VersionInfo{OSQueryVersion: osVersion, OSVersion: osqueryVersion}).Error; err != nil {
		log.Error().Err(err).Msg("failed to create versionInfo")
		return err
	}
	if err = service.db.Create(&models.InstalledApplications{AppsInstalled: appsInstalled}).Error; err != nil {
		log.Error().Err(err).Msg("failed to create installedApplications")
		return err
	}
	return nil
}
