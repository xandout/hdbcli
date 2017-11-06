package config

import (
	"encoding/json"
	"os"
)

// Configuration struct to parse the JSON config
type Configuration struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// LoadConfiguration reads a JSON file from `file_path` and returns a Configuration
func LoadConfiguration(filePath string) (*Configuration, error) {
	var config *Configuration
	configFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	decErr := jsonParser.Decode(&config)
	if decErr != nil {
		return nil, decErr
	}
	confErr := configFile.Close()
	return config, confErr
}

// NewConfiguration allows you to create a Configuration struct programmatically
func NewConfiguration(hostname string, port int, username string, password string, database string) *Configuration {
	conf := Configuration{
		Hostname: hostname,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
	return &conf
}
