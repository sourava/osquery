package models

type UserData struct {
	OSVersion      string   `json:"os_version"`
	OSQueryVersion string   `json:"os_query_version"`
	AppsInstalled  []string `json:"apps_installed"`
}
