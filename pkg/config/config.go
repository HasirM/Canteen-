package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	JWTKey     string `json:"jwt_key"`
}

func Data() Config{
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("something wrong with opening json file")
	}
	var config = Config{}
	_ = json.Unmarshal([]byte(file), &config)

	return config
}




