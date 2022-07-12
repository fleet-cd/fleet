package config

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoDB MongoDB `yaml:"mongodb"`
	Server  Server  `yaml:"server"`
}

type MongoDB struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type Server struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	BasePath string `yaml:"base_path"`
}

func Read(path string) *Config {
	log.Debug().Msg("reading fleet config")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Err(err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Err(err)
	}
	return &config
}
