package models

import "time"

type VersionInfo struct {
	ID             uint `gorm:"primaryKey"`
	OSVersion      string
	OSQueryVersion string
	CreatedAt      time.Time
}
