package osquery

import (
	"encoding/json"
	"errors"
	"github.com/sourava/secfix/external/command_executor"
	"os/exec"
)

type OsqueryClientInterface interface {
	ExecuteQuery(query string) ([]map[string]string, error)
	GetOsVersion() (string, error)
	GetOsqueryVersion() (string, error)
	GetAppsInstalled() ([]string, error)
}

type OsqueryClient struct {
	executor command_executor.CommandExecutor
}

func NewOsqueryClient(executor command_executor.CommandExecutor) *OsqueryClient {
	if executor == nil {
		executor = &command_executor.DefaultExecutor{}
	}
	return &OsqueryClient{
		executor,
	}
}

func (client *OsqueryClient) ExecuteQuery(query string) ([]map[string]string, error) {
	cmd := exec.Command("osqueryi", "--json", query)
	output, err := client.executor.CombinedOutput(cmd)
	if err != nil {
		return nil, err
	}
	var rows []map[string]string
	err = json.Unmarshal(output, &rows)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("no rows found")
	}
	return rows, nil
}

func (client *OsqueryClient) GetOsVersion() (string, error) {
	rows, err := client.ExecuteQuery("SELECT version FROM os_version LIMIT 1;")
	if err != nil {
		return "", err
	}
	if val, ok := rows[0]["version"]; ok {
		return val, nil
	} else {
		return "", errors.New("no version found")
	}
}

func (client *OsqueryClient) GetOsqueryVersion() (string, error) {
	rows, err := client.ExecuteQuery("SELECT version FROM osquery_info;")
	if err != nil {
		return "", err
	}
	if val, ok := rows[0]["version"]; ok {
		return val, nil
	} else {
		return "", errors.New("no version found")
	}
}

func (client *OsqueryClient) GetAppsInstalled() ([]string, error) {
	rows, err := client.ExecuteQuery("SELECT name FROM apps;")
	if err != nil {
		return []string{}, err
	}
	var apps []string
	for _, row := range rows {
		if val, ok := row["name"]; ok {
			apps = append(apps, val)
		}
	}
	return apps, nil
}
