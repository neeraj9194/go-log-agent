package config

import (
	// "github.com/go-yaml/yaml"
	"fmt"

	"github.com/jinzhu/configor"
)

type Config struct {
	FilePath    string
	ServiceName string
}

// func LoadConfig() Config {
// 	configLoader := NewLoader("/home/neeraj/projects/go-log-agent/config/config.yaml")

// 	config := &Config{}

// 	err := configLoader.Load(yaml.Unmarshal, config)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return *config
// }

func LoadConfig() Config {
	config := Config{}
	configor.Load(&config, "/home/neeraj/projects/go-log-agent/config/config.yaml")
	fmt.Printf("config: %#v\n", config)
	return config
}
