package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	Database Database
	Server   Server
}
type Database struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       string
}
type Server struct {
	Port int
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
		srv := Server{
			Port: viper.GetInt("server.port"),
		}
		configInstance = &Config{
			Database: db,
			Server:   srv,
		}
		fmt.Println("config initialized")
	})
	return configInstance
}
