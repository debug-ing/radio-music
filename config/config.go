package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Port   string
		Folder string
	}
}

var configPublic *Config
var once sync.Once

func LoadConfig() (config *Config) {
	once.Do(func() {
		viper.SetConfigFile("config/config.toml")
		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		if err := viper.Unmarshal(&config); err != nil {
			panic("ERROR load config file!")
		}
		configPublic = config
		log.Println("================ Loaded Configuration ================")
	})
	return
}
