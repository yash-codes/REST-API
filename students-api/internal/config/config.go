package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// env-required:"true"
type Config struct {
	Env          string `yaml:"env" env:"ENV"`
	StoaragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer   `yaml:"http_server"`
}

func MustLoad() *Config {

	// variable for configuration file path
	var configPath string

	// try to get from the environment variable
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		//check if passed through arguments/flags
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	//if err := Load(configPath, &cfg); err != nil {
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	return &cfg
}
