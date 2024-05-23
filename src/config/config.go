package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	App      AppConfig      `json:"app"`
	Database DatabaseConfig `json:"database"`
}

type AppConfig struct {
	Host      string `json:"host"`
	Port      string `json:"port"`
	SecretKey string `json:"secretkey"`
}

type DatabaseConfig struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	User   string `json:"user"`
	Pwd    string `json:"pwd"`
	Dbname string `json:"dbname"`
}

func LoadConfig(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var config Configuration
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
