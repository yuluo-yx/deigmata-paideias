package model

import (
	"reflect"
)

type Application struct {
	Port        string `yaml:"port"`
	Model       string `yaml:"model"`
	Name        string `yaml:"name"`
	ConsoleLog  bool   `yaml:"console-log"`
	LogLevel    string `yaml:"log-level"`
	LogFilePath string `yaml:"log-file-path"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type MySQL struct {
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	Databases string `yaml:"databases"`
}

type System struct {
	Mysql       MySQL       `yaml:"mysql"`
	Redis       Redis       `yaml:"redis"`
	Application Application `yaml:"application"`
}

type Config struct {
	System System `yaml:"system"`
}

var config Config

func SetConfig(tmpConfig Config) {

	if isEmpty := reflect.DeepEqual(
		reflect.ValueOf(tmpConfig),
		reflect.Zero(reflect.TypeOf(tmpConfig))); isEmpty {

		panic("Setting system configuration failed, systemInfo is null!")
	}

	config = tmpConfig
}

func GetConfig() Config {

	return config
}
