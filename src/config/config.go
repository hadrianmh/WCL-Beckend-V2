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
	Host                     string `json:"host"`
	Port                     string `json:"port"`
	JWTSecretKey             string `json:"jwtsecretkey"`
	JWTRefSecretKey          string `json:"jwtrefsecretkey"`
	JWTtokenexp              int    `json:"jwt_token_exp_minutes"`
	JWTreftokenexp           int    `json:"jwt_reftoken_exp_days"`
	DateFormat_Print         string `json:"dateformat_print"`
	DateFormat_Global        string `json:"dateformat_global"`
	DateFormat_Frontend      string `json:"dateformat_frontend"`
	DateFormat_Timestamp     string `json:"dateformat_timestamp"`
	DateFormat_Year          string `json:"dateformat_year"`
	DateFormat_Years         string `json:"dateformat_years"`
	DateFormat_Month         string `json:"dateformat_month"`
	DateFormat_Day           string `json:"dateformat_day"`
	DateFormat_MonthlyReport string `json:"dateformat_monthlyreport"`
	UploadPath               string `json:"upload_path"`
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
