package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type TestConfig struct {
	IP         string `json:"ip"`
	Port       string `json:"port"`
	RDS        string `json:"redis"`
	test_token string
	baseurl    string
	url        string
}

var testConfig TestConfig

func ReadtestConfig(ConfigFile string) {
	configHandle, err := os.Open(ConfigFile)
	if err != nil {
		log.Printlnml(log.InfoLevel, log.ServerConfigModule, "Config File Open Error:", err)
		return
	}
	defer configHandle.Close()
	configBytes, err := ioutil.ReadAll(configHandle)
	if err != nil {
		log.Printlnml(log.InfoLevel, log.ServerConfigModule, "Config File Read Error:", err)
		return
	}
	if err = json.Unmarshal(configBytes, &testConfig); err != nil {
		log.Printlnml(log.InfoLevel, log.ServerConfigModule, "Config File Content format Error:", err)
		return
	}
	return
}

type ServerConfig struct {
	RedisDBAddr         string  `json:"RedisDBAddr"`
	JWT                 bool    `json:"JWT"`
	ImageTestStart      bool    `json:"ImageTestStart"`
	NginxRoot           string  `json:"NginxRoot"`
	NginxGUIPrefix      string  `json:"NginxGUIPrefix"`
	FacePoseTestMethod  string  `json:"FacePoseTestMethod"`
	FaceClusterCntLimit int     `json:FaceClusterCntLimit`
	HistoryStreamLimit  uint32  `json:"HistoryStreamLimit"`
	DiskUsageLimit      float64 `json:"DiskUsageLimit"`
	GUIVideoDiskLimit   uint64  `json:"GUIVideoDiskLimit"`
	AIVideoDiskLimit    uint64  `json:"AIVideoDiskLimit"`
	AIPictureDiskLimit  uint64  `json:"AIPictureDiskLimit"`
}

var serverConfig ServerConfig

func ReadServerConfig(ConfigFile string) ServerConfig {
	configHandle, err := os.Open(ConfigFile)
	if err != nil {
		log.Printlnml(log.InfoLevel, log.ServerConfigModule, "Config File Open Error:", err)
	}
	defer configHandle.Close()
	configBytes, err := ioutil.ReadAll(configHandle)
	if err != nil {
		log.Printlnml(log.InfoLevel, log.ServerConfigModule, "Config File Read Error:", err)
	}
	var serverconfig ServerConfig
	if err = json.Unmarshal(configBytes, &serverconfig); err != nil {
		panic(err)
	}
	return serverconfig
}

func InitDefServerConfig() {
	// setup global nginx directories
	if serverConfig.NginxRoot == "" {
		serverConfig.NginxRoot = "/media"
	}
	if serverConfig.NginxGUIPrefix == "" {
		serverConfig.NginxGUIPrefix = "mediaassets"
	}
	if serverConfig.FaceClusterCntLimit == 0 {
		serverConfig.FaceClusterCntLimit = 200
	}

	if serverConfig.HistoryStreamLimit == 0 {
		serverConfig.HistoryStreamLimit = 10
	}
	if serverConfig.DiskUsageLimit == 0 {
		serverConfig.DiskUsageLimit = 0.95
	}
}
