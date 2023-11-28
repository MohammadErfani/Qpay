package config

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database Database
	Server   Server
	Admin    Admin
	JWT      JWT
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

type Admin struct {
	Name     string
	Username string
	Password string
	Email    string
}
type JWT struct {
	SecretKey      string
	ExpirationTime time.Duration
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
		admin := Admin{
			Name:     viper.GetString("admin.name"),
			Username: viper.GetString("admin.username"),
			Email:    viper.GetString("admin.email"),
			Password: viper.GetString("admin.password"),
		}
		jwt := JWT{
			SecretKey:      viper.GetString("auth.secret-key"),
			ExpirationTime: viper.GetDuration("auth.expiresIn"),
		}
		configInstance = &Config{
			Database: db,
			Server:   srv,
			Admin:    admin,
			JWT:      jwt,
		}

		fmt.Println("config initialized")
	})
	return configInstance
}
