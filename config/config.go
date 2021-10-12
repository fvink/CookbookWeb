package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port" envconfig:"PORT"`
		Host string `yaml:"host" envconfig:"SERVER_HOST"`
	} `yaml:"server"`
	Database struct {
		Name     string `yaml:"name" envconfig:"DB_NAME"`
		Username string `yaml:"user" envconfig:"DB_USERNAME"`
		Password string `yaml:"pass" envconfig:"DB_PASSWORD"`
		Protocol string `yaml:"protocol" envconfig:"DB_PROTOCOl"`
		Server   string `yaml:"server" envconfig:"DB_Server"`
		Port     string `yaml:"port" envconfig:"DB_PORT"`
	} `yaml:"database"`
}

func LoadConfig(filename string) (conf Config, err error) {
	err = readFile(filename, &conf)
	if err != nil {
		return
	}
	err = readEnv(&conf)
	if err != nil {
		return
	}
	return
}

func readFile(filename string, cfg *Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func readEnv(cfg *Config) error {
	err := envconfig.Process("", cfg)
	if err != nil {
		return err
	}
	return nil
}
