package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config - Структура конфигурационного файла
type Config struct {
	Asterisk struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"asterisk"`
	LogLevel string `yaml:"loglevel"`
}

// GetConfigYaml - Читаем конфиг и устанавливаем параметры приложения
func (a *MyApp) GetConfigYaml(filename string) {
	log.Info("Reading config ", filename)

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	err = yaml.Unmarshal(yamlFile, &a.config)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	a.logLevel = setLogLevel(a.config.LogLevel)
}

// setLogLevel - Устанавливаем уровень журналирования событий в приложении
func setLogLevel(confLogLevel string) log.Level {
	var result log.Level
	switch confLogLevel {
	case configLogDebug:
		result = log.DebugLevel
	case configLogInfo:
		result = log.InfoLevel
	case configLogWarn:
		result = log.WarnLevel
	case configLogError:
		result = log.ErrorLevel
	case configLogFatal:
		result = log.FatalLevel
	default:
		result = log.InfoLevel
	}

	log.Info("Application logging level: ", result)

	return result
}
