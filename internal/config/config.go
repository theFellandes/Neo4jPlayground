package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Conf struct {
	App struct {
		Port string `yaml:"port"`
	} `yaml:"app"`
	DB struct {
		Ip       string `yaml:"ip"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Port     int    `yaml:"port"`
	} `yaml:"db"`
}

const CONFIG_PATH = "internal/config/config.yaml"

func (c *Conf) GetConf() *Conf {
	yamlFile, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
