package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	ServicePort      uint
	DebugLevel       string
	PostgresHost     string
	PostgresPort     uint
	PostgresDBName   string
	PostgresUser     string
	PostgresPassword string
}

var (
	ServiceConfig Config
)

func Initialize() {
	serviceEnv := viper.New()

	serviceEnv.SetDefault("debug_level", "debug")
	serviceEnv.SetDefault("service_port", 8000)
	serviceEnv.SetDefault("postgres_host", "127.0.0.1")
	serviceEnv.SetDefault("postgres_port", 5432)
	serviceEnv.SetDefault("postgres_db_name", "testdb")
	serviceEnv.SetDefault("postgres_user", "test_user")
	serviceEnv.SetDefault("postgres_password", "test_password")

	serviceEnv.SetConfigName(".env")
	serviceEnv.SetConfigType("env")
	serviceEnv.AddConfigPath(getCwdFromExe())
	serviceEnv.AutomaticEnv()

	if !fileExists(filepath.Join(getCwdFromExe(), ".env")) {
		_, err := os.Create(filepath.Join(getCwdFromExe(), ".env"))
		if err != nil {
			log.Fatalf("[FATAL] .env doesn't exist and couldn't be created: %v", err)
		}
	}
	if err := serviceEnv.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("[FATAL] Error while reading in .env file: %v", err)
		} else {
			log.Fatalf("[FATAL]Error while parsing .env file: %v", err)
		}
	}
	setConfigFromEnv(serviceEnv)
}

func setConfigFromEnv(serviceEnv *viper.Viper) {
	ServiceConfig.DebugLevel = serviceEnv.GetString("debug_level")
	ServiceConfig.ServicePort = serviceEnv.GetUint("service_port")

	ServiceConfig.PostgresHost = serviceEnv.GetString("postgres_host")
	ServiceConfig.PostgresPort = serviceEnv.GetUint("postgres_port")
	ServiceConfig.PostgresDBName = serviceEnv.GetString("postgres_db_name")
	ServiceConfig.PostgresUser = serviceEnv.GetString("postgres_user")
	ServiceConfig.PostgresPassword = serviceEnv.GetString("postgres_password")
}

func getCwdFromExe() string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("[-] Failed to get path to current executable: %v", err)
	}
	return filepath.Dir(exe)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return !info.IsDir()
}
