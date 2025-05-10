package utils

import (
	"encoding/json"
	"os"
)

type NodeConfig struct {
	ID       string `json:"id"`
	Address  string `json:"address"`
	IsMaster bool   `json:"isMaster"`
}

type Config struct {
	Nodes []NodeConfig `json:"nodes"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

var (
	IsMaster            bool
	OtherNodes          []NodeConfig
	IsDBReady           bool
	CurrentDatabaseName string
	CurrentMaster       string
)

func InitRoles(currentAddress string, config *Config) {
	for _, node := range config.Nodes {
		if node.IsMaster {
			CurrentMaster = node.Address
		}
		if node.Address == currentAddress {
			IsMaster = node.IsMaster
		} else {
			OtherNodes = append(OtherNodes, node)
		}
	}
}
