package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Journal                 JournalConfig `json:"journal"`
	CommandDomainSocketPath string        `json:"commandSocketPath"`
	GelfUdpListeners        []string      `json:"gelfUdpListeners"`
	GelfTcpListeners        []string      `json:"gelfTcpListeners"`
	GelfHttpListeners       []string      `json:"gelfHttpListeners"`
}

type JournalConfig struct {
	DataDirPath string `json:"dataDir"`
	Type        string `json:"type"`
}

func GetConfig(filename string) (*Config, error) {

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func GetDefaultConfig() *Config {
	return &Config{Journal: JournalConfig{DataDirPath: "./data", Type: "simple"}, CommandDomainSocketPath: "./run/golfd.sock",
		GelfUdpListeners: []string{"127.0.0.1:12201"}, GelfTcpListeners: []string{"127.0.0.1:12201"},
		GelfHttpListeners: []string{"127.0.0.1:8080"}}
}
