package models

import (
	"gorm.io/datatypes"
	"time"
)

type InstalledApplications struct {
	ID            uint `gorm:"primaryKey"`
	AppsInstalled datatypes.JSON
	CreatedAt     time.Time
}
