package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Port        string `yaml:"port"`
		Host        string `yaml:"host"`
		Schema      string `yaml:"schema"`
		FullAddress string
	} `yaml:"server"`
	Redis struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		DB   int    `yaml:"databases"`
	} `yaml:"redis"`
	Logger struct {
		LogLevel string `yaml:"log-level"`
	} `yaml:"logging"`
}

func NewConfig(path string) Config {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fullAddress := fmt.Sprintf("%s%s:%s", cfg.Server.Schema, cfg.Server.Host, cfg.Server.Port)
	cfg.Server.FullAddress = fullAddress
	return cfg
}
