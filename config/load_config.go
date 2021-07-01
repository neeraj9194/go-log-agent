package config

import (
	"fmt"

	"github.com/jinzhu/configor"
)

var ServerURL string

type Config struct {
	Watchers  []Watcher
	ServerURL string
}

type Watcher struct {
	FilePath    string
	ServiceName string
}

func LoadConfig(fileName string) Config {
	config := Config{}
	configor.Load(&config, fileName)
	fmt.Printf("config: %#v\n", config)
	ServerURL = config.ServerURL
	return config
}
