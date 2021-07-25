package main

import (
	"github.com/ivahaev/amigo"
	log "github.com/sirupsen/logrus"
	"time"
)

// MyApp - Здесь все активные хэндлеры приложения
type MyApp struct {
	config   Config
	logLevel log.Level
	ami      *amigo.Amigo
}

// will be filled at build phase
var gitHash, buildTime string

func main() {
	var App MyApp

	log.Info("rt-queue-devstate version ", gitHash, " build at ", buildTime)

	// Читаем конфиг
	App.GetConfigYaml(configFileName)

	// Устанавливаем уровень журналирования событий приложения
	log.SetLevel(App.logLevel)

	// Запускаем подключение к Asterisk
	go App.RunAsteriskWorker()

	// Запускаем таймер на периодическую проверку подключения в Asterisk
	ticker := time.NewTicker(time.Duration(App.config.Asterisk.Reconnect) * time.Second)

	for {
		select {
		case <-ticker.C:
			if !App.ami.Connected() {
				App.ami.Connect()
			}
		}
	}

}
