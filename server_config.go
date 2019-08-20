package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ServerConfig global server config
type ServerConfig struct {
	LoggingReleaseMode    bool   `json:"loggingReleaseMode"`
	LoggingLevel          string `json:"loggingLevel"`
	LoggingDestination    string `json:"loggingDestination"`
	DBHostAddress         string `json:"DBHostAddress"`
	DBName                string `json:"DBName"`
	DBUsername            string `json:"DBUsername"`
	DBPassword            string `json:"DBPassword"`
	AzureStorageAccount   string `json:"azureStorageAccount"`
	AzureStorageAccessKey string `json:"azureStorageAccessKey"`
	AzureStorageContainer string `json:"azureStorageContainer"`
}

var serverConfig *ServerConfig

func readServerConfig(ConfigFile string) (*ServerConfig, error) {
	configHandle, err := os.Open(ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("Server Config File Open Error - %v", err.Error())
	}
	defer configHandle.Close()
	configBytes, err := ioutil.ReadAll(configHandle)
	if err != nil {
		return nil, fmt.Errorf("Server Config Read Error - %v", err.Error())
	}
	var serverconfig ServerConfig
	if err = json.Unmarshal(configBytes, &serverconfig); err != nil {
		return nil, fmt.Errorf("Server Config Parse Error - %v", err.Error())
	}
	return &serverconfig, nil
}

func initDefaultServerConfig(sc *ServerConfig) {
	sc.LoggingLevel = strings.ToLower(sc.LoggingLevel)
	if sc.LoggingLevel == "" {
		sc.LoggingLevel = "debug"
	}
	if sc.LoggingDestination == "" {
		sc.LoggingDestination = "stdout+file"
	}
}
