package setting

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Webset struct {
		Host          string `yaml:"host"`
		Port          int    `yaml:"port"`
		SignupEnabled bool   `yaml:"signup_enabled"`
		KwtKey        string `yaml:"kwt_key"`
	} `yaml:"webset"`
	Mysql struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		Db       int    `yaml:"db"`
	} `yaml:"redis"`
	Mongo struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"mongo"`
	Email struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		UserName string `yaml:"userName"`
		Token    string `yaml:"token"`
	} `yaml:"email"`
}

var Config *Configuration

func loadConfiguration() error {
	date, err := ioutil.ReadFile("conf/conf.yaml")
	if err != nil {
		return err
	}
	var config Configuration
	err = yaml.Unmarshal(date, &config)
	if err != nil {
		return err
	}
	Config = &config
	return nil
}

func init() {
	err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
}
