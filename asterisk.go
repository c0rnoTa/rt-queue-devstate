package main

import (
	"fmt"
	"github.com/ivahaev/amigo"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// RunAsteriskWorker - Запускает воркер Asterisk-а по чтению событий
func (a *MyApp) RunAsteriskWorker() {
	// Устанавливаем уровень журналирования событий функции
	log.SetLevel(a.logLevel)
	log.Infof("AMI Connecting to %s:%d", a.config.Asterisk.Host, a.config.Asterisk.Port)
	settings := &amigo.Settings{
		Username: a.config.Asterisk.Username,
		Port:     strconv.Itoa(a.config.Asterisk.Port),
		Password: a.config.Asterisk.Password,
		Host:     a.config.Asterisk.Host,
	}
	a.ami = amigo.New(settings)

	a.ami.Connect()

	a.ami.On("connect", func(message string) {
		log.Info("AMI connected to: ", message)
		log.Infof("Flush all custom devstates")
		_, err := a.ami.Action(map[string]string{"Action": "Command", "Command": amiFlushDevstate})
		if err != nil {
			log.Error("AMI could not flush current custom devstates: ", err)
		}
	})

	a.ami.On("error", func(message string) {
		amiConn := fmt.Sprintf("%s:%s@%s:%d", a.config.Asterisk.Username, a.config.Asterisk.Password, a.config.Asterisk.Host, a.config.Asterisk.Port)
		log.Warnf("AMI connection error [%s]: %s", amiConn, message)
	})

	err := a.ami.RegisterHandler(amiEventInUse, a.SetInuse)
	if err != nil {
		log.Error("AMI could not register handler: ", err)
	}

	err = a.ami.RegisterHandler(amiEventNotInUse, a.SetNotinuse)
	if err != nil {
		log.Error("AMI could not register handler: ", err)
	}

}

// SetInuse - Устанавливает Custom dev state в IN_USE
func (a *MyApp) SetInuse(m map[string]string) {
	log.SetLevel(a.logLevel)
	log.Debugf("AMI event received: %v\n", m)
	fields, err := getFields(m, amiFieldMember)
	if err != nil {
		log.Error("AMI Error in event handler: ", err)
		return
	}
	log.Infof("AMI action here for member %s", fields[amiFieldMember])
	_, err = a.ami.Action(map[string]string{"Action": "Command", "Command": fmt.Sprintf(amiCommand, fields[amiFieldMember], devInuse)})
}

// SetNotinuse - Устанавливает Custom dev state в NOT_INUSE
func (a *MyApp) SetNotinuse(m map[string]string) {
	log.SetLevel(a.logLevel)
	log.Debugf("AMI event received: %v\n", m)
	fields, err := getFields(m, amiFieldMember)
	if err != nil {
		log.Error("AMI Error in event handler: ", err)
		return
	}
	log.Infof("AMI action here for member %s", fields[amiFieldMember])
	_, err = a.ami.Action(map[string]string{"Action": "Command", "Command": fmt.Sprintf(amiCommand, fields[amiFieldMember], devNotInuse)})
}

// getFields - Получает значения параметров из события
func getFields(m map[string]string, fields ...string) (map[string]string, error) {
	values := make(map[string]string)
	for _, field := range fields {
		value, ok := m[field]
		if !ok {
			log.WithFields(log.Fields{
				"map": m,
			}).Error("Invalid params map")
			// TODO FIX error handling here
			return nil, nil
		}
		values[field] = value
	}
	return values, nil
}
