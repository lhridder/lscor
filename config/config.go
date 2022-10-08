package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

var GlobalConfig Config

type Discord struct {
	Token   string `yaml:"token"`
	Guild   string `yaml:"guild"`
	Message string `yaml:"message"`
	Channel string `yaml:"channel"`
}

type Corero struct {
	URL  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type Config struct {
	Discord  Discord
	Corero   Corero
	Interval struct {
		Notification string `yaml:"notification"`
		Embed        string `yaml:"embed"`
	}
	Listen string `yaml:"listen"`
	Debug  bool   `yaml:"debug"`
}

func Load() error {
	log.Println("Loading config.yml")
	ymlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(ymlFile, &GlobalConfig)
	if err != nil {
		return err
	}
	return nil
}
