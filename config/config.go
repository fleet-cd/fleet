package config

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoDB     MongoDB     `yaml:"mongodb"`
	Server      Server      `yaml:"server"`
	Git         Git         `yaml:"git"`
	Artifactory Artifactory `yaml:"artifactory"`
	Kubernetes  Kubernetes  `yaml:"kubernetes"`
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

type Kubernetes struct {
	ConfigPath string `yaml:"config_path"`
}

type Git struct {
	GithubAccessToken string `yaml:"github_access_token"`
}

type Artifactory struct {
	AccessToken string `yaml:"access_token"`
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
