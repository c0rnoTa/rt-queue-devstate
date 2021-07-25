## Asterisk devstate aid for realtime chan_sip members

### Описание проблемы

**Предпосылки:** 

* Вы используете chan_sip для подключения своих абонентов.
* Абоненты хранятся в базе данных через технологию [REALTIME](https://wiki.asterisk.org/wiki/display/AST/Realtime+Database+Configuration#RealtimeDatabaseConfiguration-RealtimeSIPfriends)
* Абоненты являются операторами очереди (members)
* В очереди установлено ограничение вызовов на операторов, которые находятся в разговоре (`ringinuse=no`)  
* Ваш оператор находится в разговоре с вызывающим из очереди. Его devstate находится в состоянии `in use`

**Ситуация:** 

1. Вы делаете `sip reload`, оператор всё еще находится в разговоре. devstate находится в состоянии `in use`. 
2. Абонентское устройство перерегистрируется и devstate меняется на `Not in use`, хотя абонент остаётся в разговоре.
3. На оператора начинают поступать параллельные звонки, хотя `ringinuse=no`.

### Решение

**Концепция:**

* Хранить состояние доступности оператора в очереди dev_state в каком-нибудь другом месте, например, в [custom devstate](https://wiki.asterisk.org/wiki/display/AST/Device+State#DeviceState-CustomDeviceStates).
* Обновлять custom dev state в зависимости от событий по обработке звонков оператором. Т.е. опираться на события для того,
чтобы понимать реальное состояние оператора.
  
**Реализация:**

1. При определении оператора очереди (member) указывать `state_interface` как `hint`
2. Сделать совмещенный `hint`, как объединение реального состояния абонента (от канала), так и custom. 
При этом, если одно из состояний будет `In use`, то `hint` будет возвращать `In use`.
3. Использовать сервис из этого репозитория, который будет слушать события в AMI и менять `Custom` dev_state.
Сервис ориентируется на события AMI `AgentConnect` и `AgentComplete`. 

### Сборка

```shell
go build -ldflags "-X main.buildTime=`date +%Y-%m-%d:%H:%M:%S` -X main.gitHash=`git rev-parse --short HEAD`"
```

### Установка

1. В `/etc/asterisk/manager.conf` добавляем 
```
[rt-queue-devstate]
secret = password
deny = 0.0.0.0/0.0.0.0
permit = 127.0.0.1/255.255.255.255
read = system,call,log,verbose,agent,user,config,dtmf,reporting,cdr,dialplan
write = system,call,agent,user,config,command,reporting,originate,message
eventfilter=Event: AgentComplete
eventfilter=Event: AgentConnect
```

2. Бинарь размещаем в одной директории с конфигом в `/opt/fibex/rt-queue-devstate`
3. Systemd юнит добавляем в `/etc/systemd/system` и делаем `systemctl daemon-reload`