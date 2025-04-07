package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	userDataHandler "github.com/sourava/secfix/internal/handler/user_data"
	"github.com/sourava/secfix/internal/models"
	userDataService "github.com/sourava/secfix/internal/service/user_data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func main() {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to database")
		return
	}
	fmt.Println("Connected to database", db)

	err = db.AutoMigrate(&models.VersionInfo{})
	if err != nil {
	}
	err = db.AutoMigrate(&models.InstalledApplications{})
	if err != nil {
	}

	userDataServiceObj := userDataService.NewUserDataService(db)
	userDataHandlerObj := userDataHandler.NewUserDataHandler(userDataServiceObj)

	r := gin.Default()
	r.GET("/latest_data", userDataHandlerObj.GetLatestUserData)
	err = r.Run()
	if err != nil {
		return
	}
}
