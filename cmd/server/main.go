package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/sourava/secfix/external/osquery"
	userDataHandler "github.com/sourava/secfix/internal/handler/user_data"
	"github.com/sourava/secfix/internal/models"
	userDataService "github.com/sourava/secfix/internal/service/user_data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

func addUserData(userDataService *userDataService.UserDataService) {
	log.Info().Msg("Adding User Data at " + time.Now().String())
	err := userDataService.AddLatestUserData()
	if err != nil {
		log.Error().Err(err).Msg("error adding latest user data")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("error loading .env file")
	}

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("error connecting to the database")
		return
	}
	log.Info().Msg("connected to the database")

	err = db.AutoMigrate(&models.VersionInfo{})
	if err != nil {
		log.Error().Err(err).Msg("error auto-migrating version info")
	}
	err = db.AutoMigrate(&models.InstalledApplications{})
	if err != nil {
		log.Error().Err(err).Msg("error auto-migrating installed applications")
	}

	osqueryClient := osquery.NewOsqueryClient(nil)
	userDataServiceObj := userDataService.NewUserDataService(db, osqueryClient)
	userDataHandlerObj := userDataHandler.NewUserDataHandler(userDataServiceObj)

	// Updating data on app startup
	addUserData(userDataServiceObj)
	// Updating data every minute
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				addUserData(userDataServiceObj)
			}
		}
	}()

	r := gin.Default()
	r.GET("/latest_data", userDataHandlerObj.GetLatestUserData)
	err = r.Run()
	if err != nil {
		log.Error().Err(err).Msg("error starting server")
		return
	}
}
