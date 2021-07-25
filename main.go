package main

import (
	"github.com/ivahaev/amigo"
	log "github.com/sirupsen/logrus"
)

// MyApp - Здесь все активные хэндлеры приложения
type MyApp struct {
	config   Config
	logLevel log.Level
	ami      *amigo.Amigo
}

// will be filled at buid phase
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

	// TODO Сюда можно добавить проверку статуса подключений и их восстановление в цикле
	ch := make(chan bool)
	<-ch
}
