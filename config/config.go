package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Database Database
}
type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       string
}

var (
	configInstance *Config
	once           sync.Once
)

func InitConfig(configPath string) *Config {
	once.Do(func() {
		viper.SetConfigFile(configPath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to open config file: %v", err))
		}
		db := Database{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			DB:       viper.GetString("database.db"),
		}
		configInstance = &Config{
			Database: db,
		}
		fmt.Println("config initialized")
	})
	return configInstance
}
