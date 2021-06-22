package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

type Config struct {
	FilePath    string
	ServiceName string
	ServerURL   string
}

func LoadConfig(fileName string) Config {
	config := Config{}
	configor.Load(&config, fileName)
	fmt.Printf("config: %#v\n", config)
	return config
}
