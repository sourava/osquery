package user_data

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sourava/secfix/external/osquery"
	"github.com/sourava/secfix/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestGetLatestUserData_ShouldReturnError_WhenDBReturnsErrorWhileQuerying(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dial := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dial, &gorm.Config{})
	mockOsqueryClient := new(osquery.MockOsqueryClient)
	service := NewUserDataService(db, mockOsqueryClient)
	mock.ExpectQuery(`SELECT`).WillReturnError(errors.New("some error"))

	data, err := service.GetLatestUserData()

	assert.Nil(t, data)
	assert.NotNil(t, err)
}

func TestGetLatestUserData_ShouldReturnError_WhenAppsInstalledAreNotStoredInCorrectFormat(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dial := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dial, &gorm.Config{})
	mockOsqueryClient := new(osquery.MockOsqueryClient)
	service := NewUserDataService(db, mockOsqueryClient)
	versionInfoRows := sqlmock.NewRows([]string{"os_version", "os_query_version"}).AddRow("1.0.0", "1.0.0")
	mock.ExpectQuery(`SELECT`).WillReturnRows(versionInfoRows)
	installedAppsRows := sqlmock.NewRows([]string{"apps_installed"}).AddRow("{}")
	mock.ExpectQuery(`SELECT`).WillReturnRows(installedAppsRows)

	data, err := service.GetLatestUserData()

	assert.Nil(t, data)
	assert.NotNil(t, err)
}

func TestGetLatestUserData_ShouldReturnData(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dial := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dial, &gorm.Config{})
	mockOsqueryClient := new(osquery.MockOsqueryClient)
	service := NewUserDataService(db, mockOsqueryClient)
	versionInfoRows := sqlmock.NewRows([]string{"os_version", "os_query_version"}).AddRow("1.0.0", "1.0.0")
	mock.ExpectQuery(`SELECT`).WillReturnRows(versionInfoRows)
	installedAppsRows := sqlmock.NewRows([]string{"apps_installed"}).AddRow("[\"App1\"]")
	mock.ExpectQuery(`SELECT`).WillReturnRows(installedAppsRows)

	data, err := service.GetLatestUserData()

	assert.Equal(t, &models.UserData{
		OSVersion:      "1.0.0",
		OSQueryVersion: "1.0.0",
		AppsInstalled:  []string{"App1"},
	}, data)
	assert.Nil(t, err)
}

func TestAddLatestUserData_ShouldReturnError_WhenDBReturnsErrorWhileInserting(t *testing.T) {
	mockDb, mock, _ := sqlmock.New()
	dial := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dial, &gorm.Config{})
	mockOsqueryClient := new(osquery.MockOsqueryClient)
	mockOsqueryClient.On("GetOsVersion").Return("1.0.0", nil)
	mockOsqueryClient.On("GetOsqueryVersion").Return("1.0.0", nil)
	mockOsqueryClient.On("GetAppsInstalled").Return([]string{"App1", "App2"}, nil)

	service := NewUserDataService(db, mockOsqueryClient)
	mock.ExpectQuery(`INSERT`).WillReturnError(errors.New("some error"))

	err := service.AddLatestUserData()

	assert.NotNil(t, err)
}
