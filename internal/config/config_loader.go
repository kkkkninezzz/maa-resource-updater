package config

import (
	"encoding/json"
	"log"
	"os"
)

func LoadMaaResouceUpdaterConfig(configPath string) MaaResourceUpdaterConfig {

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
		panic(err)
	}

	// 解析 JSON 数据
	var config MaaResourceUpdaterConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		panic(err)
	}

	return config
}
