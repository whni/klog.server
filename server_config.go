package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// ServerConfig global server config
type ServerConfig struct {
	LoggingReleaseMode bool   `json:"loggingReleaseMode"`
	LoggingLevel       string `json:"loggingLevel"`
}

var serverConfig *ServerConfig

func readServerConfig(ConfigFile string) *ServerConfig {
	configHandle, err := os.Open(ConfigFile)
	if err != nil {
		log.Println(log.InfoLevel, "Config File Open Error:", err)
	}
	defer configHandle.Close()
	configBytes, err := ioutil.ReadAll(configHandle)
	if err != nil {
		log.Println(log.InfoLevel, "Config File Read Error:", err)
	}
	var serverconfig ServerConfig
	if err = json.Unmarshal(configBytes, &serverconfig); err != nil {
		panic(err)
	}
	return &serverconfig
}

func initDefaultServerConfig(sc *ServerConfig) {
	sc.LoggingLevel = strings.ToLower(sc.LoggingLevel)
	if sc.LoggingLevel == "" {
		sc.LoggingLevel = "debug"
	}
}
